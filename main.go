package main

import (
	"bytes"
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"net"
	"net/http"
	// "reflect"
	"sort"
	"time"

	"github.com/Ullaakut/nmap"
)

func scan_subnet() []net.IP {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Equivalent to `/usr/local/bin/nmap -sP 172.16.0.0/24, with a 1 minute timeout.
	scanner, err := nmap.NewScanner(
		nmap.WithTargets("172.16.0.0/24"),
		nmap.WithPingScan(),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, err := scanner.Run()
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	realIPs := make([]net.IP, 0, 255)
	for _, host := range result.Hosts {
		realIPs = append(realIPs, net.ParseIP(host.Addresses[0].String()))
	}

	sort.Slice(realIPs, func(i, j int) bool {
		return bytes.Compare(realIPs[i], realIPs[j]) < 0
	})
	return realIPs
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(`{"author": "zTeeed", 
			 "follow_me": "https://github.com/zteeed", 
			 "paths": {"0": "/ip", "1": "/ip"}}`))
}

func ip(w http.ResponseWriter, r *http.Request) {
	result := scan_subnet()
	data := map[int]string{}
	for k, v := range result {
		// fmt.Printf("key: %d | value: %s\n", k, v)
		data[k] = v.String()
	}
	// fmt.Println(data)
	dict, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(dict))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ip", ip)
	http.ListenAndServe(":5001", nil)
}
