package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"regexp"
	"net/http"
	"encoding/json"
	"../leaseparsers"
)

type LeaseData struct {
	Data []leaseparsers.Lease
}

func getJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func dumpLeases(leases []leaseparsers.Lease, filter *regexp.Regexp) {
	fmt.Printf("%30s\t%s\t%s\t%s\n", "Hostname", "IP", "Mac", "MacVendor")

	for _, l := range leases {
		if filter != nil {
			if (filter.MatchString(l.Hostname) == false) {
				continue;
			}
		}
		fmt.Printf("%30s\t%s\t%s\t%s\n", l.Hostname, l.Ip, l.Mac, l.MacVendor)
	}
}

func usage() {
	fmt.Printf("Usage: %s [hostname-filter (Regexp)]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	uriStr := flag.String("uri", "http://localhost:8080/leases", "URI to fetch data from")
	flag.Parse()

	var filter *regexp.Regexp

	if flag.NArg() > 0 {
		filter = regexp.MustCompile(flag.Arg(0))
	}

	leases := LeaseData{}
	getJson(*uriStr, &leases)
	dumpLeases(leases.Data, filter)
}
