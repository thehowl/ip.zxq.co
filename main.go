package main

import (
	"encoding/json"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

func main() {
	// Initialize the database.
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	// Whatever we do, we should always close the database.
	defer db.Close()
	var i string
	// Tell the user to type in an IP.
	fmt.Print("Type in an IP... ")
	// Get the IP from stdin
	_, err = fmt.Scanf("%s", &i)
	if err != nil {
		log.Fatal(err)
	}
	// Parse the ip.
	ip := net.ParseIP(i)
	if ip == nil {
		fmt.Println("That's not an IP, you dumb fuck!")
		return
	}
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	// Print out the ISO code of the country.
	output, _ := json.Marshal(map[string]interface{}{"ip": ip, "country": record.Country.IsoCode})
	fmt.Printf("%s\n", output)
}
