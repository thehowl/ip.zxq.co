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
func IPToResponse(i string) (string, string) {
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
	// Print out the ISO code of the country.
	output, _ := json.Marshal(map[string]interface{}{"ip": ip, "country": record.Country.IsoCode, "country_full": record.Country.Names["en"], "city": record.City.Names["en"], "region": sd, "continent": record.Continent.Code, "continent_full": record.Continent.Names["en"], "postal": record.Postal.Code, "loc": fmt.Sprintf("%.4f,%.4f", record.Location.Latitude, record.Location.Longitude)})
	return string(output[:]), "application/json"
}
func HTTPRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[rq]", r.Method, r.URL.Path)
	o, contentType := IPToResponse(strings.Split(r.URL.Path, "/")[1])
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	fmt.Fprint(w, o)
}
