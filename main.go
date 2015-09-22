package main

import (
	"encoding/json"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var db *geoip2.Reader

func main() {
	// Initialize the database.
	var err error
	db, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", HTTPRequestHandler)
	log.Println("Server listening!")
	http.ListenAndServe(":8080", nil)
}
func HTTPRequestHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestedThings := strings.Split(r.URL.Path, "/")
	var IPAddress string
	var Which string
	// How in the world the user would manage to even send a request to
	// something without even having Path = "/"?
	// I... have no idea. But I'm paranoid. So let's just do it anyway.
	if len(requestedThings) < 2 {
		IPAddress = ""
	} else {
		IPAddress = requestedThings[1]
	}
	// In case the user didn't write a specific index, let's specify it for
	// them.
	if len(requestedThings) < 3 {
		Which = ""
	} else {
		Which = requestedThings[2]
	}
	queryParamsRaw, _ := url.ParseQuery(r.URL.RawQuery)
	queryParams := SimplifyQueryMap(queryParamsRaw)
	queryParams = AppendDefaultIfNotSet(queryParams, "callback", "#none#")
	queryParams = AppendDefaultIfNotSet(queryParams, "pretty", "0")
	o, contentType := IPToResponse(IPAddress, Which, queryParams)
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	fmt.Fprint(w, o)
	log.Printf("[rq] %s %s %dns", r.Method, r.URL.Path, time.Since(start).Nanoseconds())
}
func AppendDefaultIfNotSet(sl map[string]string, k string, dv string) map[string]string {
	if _, ok := sl[k]; !ok {
		sl[k] = dv
	}
	return sl
}
func SimplifyQueryMap(sl url.Values) map[string]string {
	var ret map[string]string = map[string]string{}
	for k, v := range sl {
		// We're getting only the last element, because we take as granted that
		// what the use actually means is the last element, if he has provided
		// multiple values for the same key.
		if len(v) > 0 {
			ret[k] = v[len(v)-1]
		}
	}
	return ret
}
func IPToResponse(i string, specific string, params map[string]string) (string, string) {
	// Parse the ip.
	ip := net.ParseIP(i)
	if ip == nil {
		return "Please provide a valid IP address", "text/html"
	}
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	var sd string
	if record.Subdivisions != nil {
		sd = record.Subdivisions[0].Names["en"]
	}

	data := map[string]string{}
	data["ip"] = ip.String()
	data["country"] = record.Country.IsoCode
	data["country_full"] = record.Country.Names["en"]
	data["city"] = record.City.Names["en"]
	data["region"] = sd
	data["continent"] = record.Continent.Code
	data["continent_full"] = record.Continent.Names["en"]
	data["postal"] = record.Postal.Code
	data["loc"] = fmt.Sprintf("%.4f,%.4f", record.Location.Latitude, record.Location.Longitude)

	if specific == "" || specific == "json" || specific == "geo" {
		var bytes_output []byte
		if params["pretty"] == "1" {
			bytes_output, _ = json.MarshalIndent(data, "", "  ")
		} else {
			bytes_output, _ = json.Marshal(data)
		}
		return string(bytes_output[:]), "application/json"
	} else if val, ok := data[specific]; ok {
		return val, "text/html"
	} else {
		return "undefined", "text/html"
	}
}
