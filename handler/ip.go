package handler

import (
	"net/http"
	"context"
	"encoding/json"
	"github.com/HackademINT/nmap.priv.hackademint.org/gateway"
	"net"
	"sort"
	"bytes"
)

type ipHandler struct {
	gateway gateway.NmapGateway
}

func NewIPHandler(nmapGateway gateway.NmapGateway) (http.Handler, error) {
	return &ipHandler{
		gateway: nmapGateway,
	}, nil
}

func (h *ipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	ipAddresses, err := h.gateway.ScanSubnet(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapIPToStrings(ipAddresses))
}

func mapIPToStrings(ipAddresses []net.IP) []string {
	ipAddressesSorted := sortIPAddresses(ipAddresses)

	ipAddressesStr := make([]string, 0, len(ipAddressesSorted))
	for _, v := range ipAddressesSorted {
		ipAddressesStr = append(ipAddressesStr, v.String())
	}

	return ipAddressesStr

}
func sortIPAddresses(ipAddresses []net.IP) []net.IP {
	result := make([]net.IP, 0, len(ipAddresses))

	for _, ip := range ipAddresses {
		result = append(result, ip)
	}

	sort.Slice(result, func(i, j int) bool {
		return bytes.Compare(result[i], result[j]) < 0
	})
	return result
}
