package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
)

// The GeoIP database containing data on what IP match to what city/country blah
// blah.
var db *geoip2.Reader
var currFilename = time.Now().Format("2006-01") + ".mmdb"
var dbMtx = new(sync.RWMutex)

const dbURL = "https://download.db-ip.com/free/dbip-city-lite-%s.mmdb.gz"

func doUpdate() {
	fmt.Fprintln(os.Stderr, "Fetching updates...")
	currMonth := time.Now().Format("2006-01")
	if currFilename == currMonth+".mmdb" {
		fmt.Fprintln(os.Stderr, "Version is latest, not fetching")
		return
	}
	resp, err := http.Get(fmt.Sprintf(dbURL, currMonth))
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while attempting to fetch the updated DB: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Status code is not 200 for new database (%d), probably need to wait...\n", resp.StatusCode)
		return
	}
	// status code is 200, download file
	dst, err := os.Create(currMonth + ".mmdb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return
	}
	failed := true
	defer func() {
		dst.Close()
		if failed {
			os.Remove(dst.Name())
		}
	}()
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating gzip decoder: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stderr, "Copying new database...")
	_, err = io.Copy(dst, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error copying file: %v\n", err)
		return
	}
	newDB, err := geoip2.Open(currMonth + ".mmdb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening new db: %v\n", err)
		return
	}

	// actual update
	old := currFilename
	dbMtx.Lock()
	currFilename = currMonth + ".mmdb"
	if db != nil {
		db.Close()
	}
	db = newDB
	dbMtx.Unlock()
	if err := os.Remove(old); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error removing old file: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "GeoIP database updated to %s\n", currMonth)
}

func updater() {
	for range time.Tick(time.Hour * 24 * 7) {
		doUpdate()
	}
}

func main() {
	// Initialize the database.
	var err error
	db, err = geoip2.Open(currFilename)
	if err != nil {
		if os.IsNotExist(err) {
			currFilename = ""
			doUpdate()
			if db == nil {
				os.Exit(1)
			}
		} else {
			log.Fatal(err)
		}
	}
	go updater()

	// Get the HTTP server rollin'
	log.Println("Server listening!")
	http.ListenAndServe(":80", http.HandlerFunc(handler))
}

var invalidIPBytes = []byte("Please provide a valid IP address.")

type dataStruct struct {
	IP            string `json:"ip"`
	City          string `json:"city"`
	Region        string `json:"region"`
	Country       string `json:"country"`
	CountryFull   string `json:"country_full"`
	Continent     string `json:"continent"`
	ContinentFull string `json:"continent_full"`
	Loc           string `json:"loc"`
	Postal        string `json:"postal"`
}

var nameToField = map[string]func(dataStruct) string{
	"ip":             func(d dataStruct) string { return d.IP },
	"city":           func(d dataStruct) string { return d.City },
	"region":         func(d dataStruct) string { return d.Region },
	"country":        func(d dataStruct) string { return d.Country },
	"country_full":   func(d dataStruct) string { return d.CountryFull },
	"continent":      func(d dataStruct) string { return d.Continent },
	"continent_full": func(d dataStruct) string { return d.ContinentFull },
	"loc":            func(d dataStruct) string { return d.Loc },
	"postal":         func(d dataStruct) string { return d.Postal },
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get the current time, so that we can then calculate the execution time.
	start := time.Now()

	// Log how much time it took to respond to the request, when we're done.
	defer func() {
		log.Printf(
			"[rq] %s %s %s",
			r.Method,
			r.URL.Path,
			time.Since(start).String())
	}()

	// Separate two strings when there is a / in the URL requested.
	requestedThings := strings.Split(r.URL.Path, "/")

	var IPAddress string
	var Which string
	switch len(requestedThings) {
	case 3:
		Which = requestedThings[2]
		fallthrough
	case 2:
		IPAddress = requestedThings[1]
	}

	// Set the requested IP to the user's request request IP, if we got no address.
	if IPAddress == "" || IPAddress == "self" {
		// The request is most likely being done through a reverse proxy.
		if realIP, ok := r.Header["X-Real-Ip"]; ok && len(r.Header["X-Real-Ip"]) > 0 {
			IPAddress = realIP[0]
		} else {
			// Get the real actual request IP without the trolls
			IPAddress = UnfuckRequestIP(r.RemoteAddr)
		}
	}

	ip := net.ParseIP(IPAddress)
	if ip == nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(invalidIPBytes)
		return
	}

	// Query the maxmind database for that IP address.
	dbMtx.RLock()
	record, err := db.City(ip)
	dbMtx.RUnlock()
	if err != nil {
		log.Fatal(err)
	}

	// String containing the region/subdivision of the IP. (E.g.: Scotland, or
	// California).
	var sd string
	// If there are subdivisions for this IP, set sd as the first element in the
	// array's name.
	if record.Subdivisions != nil {
		sd = record.Subdivisions[0].Names["en"]
	}

	// Fill up the data array with the geoip data.
	d := dataStruct{
		IP:            ip.String(),
		Country:       record.Country.IsoCode,
		CountryFull:   record.Country.Names["en"],
		City:          record.City.Names["en"],
		Region:        sd,
		Continent:     record.Continent.Code,
		ContinentFull: record.Continent.Names["en"],
		Postal:        record.Postal.Code,
		Loc:           fmt.Sprintf("%.4f,%.4f", record.Location.Latitude, record.Location.Longitude),
	}

	// Since we don't have HTML output, nor other data from geo data,
	// everything is the same if you do /8.8.8.8, /8.8.8.8/json or /8.8.8.8/geo.
	if Which == "" || Which == "json" || Which == "geo" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		callback := r.URL.Query().Get("callback")
		enableJSONP := callback != "" && len(callback) < 2000 && callbackJSONP.MatchString(callback)
		if enableJSONP {
			_, err = w.Write([]byte("/**/ typeof " + callback + " === 'function' " +
				"&& " + callback + "("))
			if err != nil {
				return
			}
		}
		enc := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "1" {
			enc.SetIndent("", "  ")
		}
		enc.Encode(d)
		if enableJSONP {
			w.Write([]byte(");"))
		}
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if val := nameToField[Which]; val != nil {
			w.Write([]byte(val(d)))
		} else {
			w.Write([]byte("undefined"))
		}
	}
}

// Very restrictive, but this way it shouldn't completely fuck up.
var callbackJSONP = regexp.MustCompile(`^[a-zA-Z_\$][a-zA-Z0-9_\$]*$`)

// Remove from the IP eventual [ or ], and remove the port part of the IP.
func UnfuckRequestIP(ip string) string {
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	ss := strings.Split(ip, ":")
	ip = strings.Join(ss[:len(ss)-1], ":")
	return ip
}
