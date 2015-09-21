package main

import (
	"encoding/json"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"strings"
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
	fmt.Println("Server listening!")
	http.ListenAndServe(":8080", nil)
}
func IPToResponse(i string, specific string) (string, string) {
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
	data := map[string]string{"ip": ip.String(), "country": record.Country.IsoCode, "country_full": record.Country.Names["en"], "city": record.City.Names["en"], "region": sd, "continent": record.Continent.Code, "continent_full": record.Continent.Names["en"], "postal": record.Postal.Code, "loc": fmt.Sprintf("%.4f,%.4f", record.Location.Latitude, record.Location.Longitude)}
	if specific == "" || specific == "json" || specific == "geo" {
		bytes_output, _ := json.Marshal(data)
		return string(bytes_output[:]), "application/json"
	} else if val, ok := data[specific]; ok {
		return val, "text/html"
	} else {
		return "undefined", "text/html"
	}
}
func HTTPRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[rq]", r.Method, r.URL.Path)
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
	o, contentType := IPToResponse(IPAddress, Which)
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	fmt.Fprint(w, o)
}
