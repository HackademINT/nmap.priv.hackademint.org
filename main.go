package main

import (
	"github.com/HackademINT/nmap.priv.hackademint.org/gateway"
	"github.com/HackademINT/nmap.priv.hackademint.org/handler"
	"log"
	"net/http"
)

const _targets = "172.16.0.0/24"

func main() {
	indexHandler, err := handler.NewIndexHandler()
	if err != nil {
		log.Fatalf("could not create index handler: %v", indexHandler)
	}

	nmapGateway := gateway.NewNmapGateway(_targets)
	ipHandler, err := handler.NewIPHandler(nmapGateway)
	if err != nil {
		log.Fatalf("could not create IP handler: %v", indexHandler)
	}

	http.Handle("/", indexHandler)
	http.Handle("/ip", ipHandler)

	http.ListenAndServe(":5001", nil)
}
