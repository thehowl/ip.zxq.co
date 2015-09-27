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

// The GeoIP database containing data on what IP match to what city/country blah
// blah.
var db *geoip2.Reader

func main() {
	// Initialize the database.
	var err error
	db, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	// Get the HTTP server rollin'
	http.HandleFunc("/", HTTPRequestHandler)
	log.Println("Server listening!")
	http.ListenAndServe(":61430", nil)
}

// Standard request handler if there's no static file to be served.
func HTTPRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current time, so that we can then calculate the execution time.
	start := time.Now()

	var requestIP string
	// The request is most likely being done through a reverse proxy.
	if realIP, ok := r.Header["X-Real-Ip"]; ok && len(r.Header["X-Real-Ip"]) > 0 {
		requestIP = realIP[0]
	} else {
		// Get the real actual request IP without the trolls
		requestIP = UnfuckRequestIP(r.RemoteAddr)
	}

	// Log how much time it took to respond to the request, when we're done.
	defer log.Printf(
		"[rq] %s %s %s %dns",
		requestIP,
		r.Method,
		r.URL.Path,
		time.Since(start).Nanoseconds())

	// Index, redirect to github.com page.
	if r.URL.Path == "/" {
		http.Redirect(w, r, "https://github.com/TheHowl/ip.zxq.co/blob/master/README.md", 301)
		return
	}

	// Separate two strings when there is a / in the URL requested.
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

	// Query parameters array making
	queryParamsRaw, _ := url.ParseQuery(r.URL.RawQuery)
	queryParams := SimplifyQueryMap(queryParamsRaw)
	queryParams = AppendDefaultIfNotSet(queryParams, "callback", "#none#")
	queryParams = AppendDefaultIfNotSet(queryParams, "pretty", "0")

	// Get the geodata of the requested IP.
	o, contentType := IPToResponse(IPAddress, Which, queryParams)

	// Set the content type as the one given by IPToResponse.
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	// Write the data out to the response.
	fmt.Fprint(w, o)
}

// Appends a default value to a map only if the key, defined as k, doesn't
// already exist in the array.
func AppendDefaultIfNotSet(sl map[string]string, k string, dv string) map[string]string {
	if _, ok := sl[k]; !ok {
		sl[k] = dv
	}
	return sl
}

// url.ParseQuery returns a map containing as a value a slice with often just
// one value. We're fixing that.
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

// Remove from the IP eventual [ or ], and remove the port part of the IP.
func UnfuckRequestIP(ip string) string {
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	ss := strings.Split(ip, ":")
	ip = strings.Join(ss[:len(ss)-1], ":")
	return ip
}

// Turn the IP into a JSON string containing geodata.
//
// * i: the raw IP string.
// * specific: the specific value to get from the geodata array. Default is ""
// * params: Set callback in the map to a non-"#none#" value to use it as a
//   JSONP callback. Set "pretty" to 1 if you want a 2-space indented JSON
//   output.
func IPToResponse(i string, specific string, params map[string]string) (string, string) {
	ip := net.ParseIP(i)
	if ip == nil {
		return "Please provide a valid IP address", "text/html"
	}

	// Query the maxmind database for that IP address.
	record, err := db.City(ip)
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

	// Create a new instance of all the data to be returned to the user.
	data := map[string]string{}
	// Fill up the data array with the geoip data.
	data["ip"] = ip.String()
	data["country"] = record.Country.IsoCode
	data["country_full"] = record.Country.Names["en"]
	data["city"] = record.City.Names["en"]
	data["region"] = sd
	data["continent"] = record.Continent.Code
	data["continent_full"] = record.Continent.Names["en"]
	data["postal"] = record.Postal.Code
	// precision of latitude/longitude is up to 4 decimal places (even on
	// ipinfo.io).
	data["loc"] = fmt.Sprintf("%.4f,%.4f", record.Location.Latitude, record.Location.Longitude)

	// Since we don't have HTML output, nor other data from geo data,
	// everything is the same if you do /8.8.8.8, /8.8.8.8/json or /8.8.8.8/geo.
	if specific == "" || specific == "json" || specific == "geo" {
		var bytes_output []byte
		if params["pretty"] == "1" {
			bytes_output, _ = json.MarshalIndent(data, "", "  ")
		} else {
			bytes_output, _ = json.Marshal(data)
		}
		return JSONPify(params["callback"], string(bytes_output[:])),
			"application/json"
	} else if val, ok := data[specific]; ok {
		// If we got a specific value for what the user requested, return only
		// that specific value.
		return val, "text/html"
	} else {
		// We got nothing to show to the user.
		return "undefined", "text/html"
	}
}

// Wraps wrapData into a JSONP callback, if the callback name is valid.
func JSONPify(callback string, wrapData string) string {
	// If you have a callback name longer than 2000 characters, I gotta say, you
	// really should learn to minify your javascript code!
	if callback != "#none#" && callback != "" && len(callback) < 2000 {
		// In case you're wondering, yes, there is a reason for the empty
		// comment! http://stackoverflow.com/a/16048976/5328069
		wrapData = fmt.Sprintf("/**/ typeof %s === 'function' "+
			"&& %s(%s);", callback, callback, wrapData)
	}
	return wrapData
}
