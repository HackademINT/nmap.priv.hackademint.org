package gateway

import (
	"context"
	"net"
	"time"
	"fmt"
	"github.com/Ullaakut/nmap"
	"golang.org/x/sync/singleflight"
)

const _cacheDuration = 1 * time.Minute

type NmapGateway interface {
	ScanSubnet(ctx context.Context) ([]net.IP, error)
}

type nmapGateway struct {
	targets string
	group   singleflight.Group
	cache   []net.IP
}

func NewNmapGateway(targets string) NmapGateway {
	return &nmapGateway{
		targets: targets,
		group:   singleflight.Group{},
	}
}

func (g *nmapGateway) ScanSubnet(ctx context.Context) ([]net.IP, error) {
	if g.cache != nil {
		return g.cache, nil
	}

	result, err, _ := g.group.Do("scanSubnet", func() (interface{}, error) {
		return g.refreshCache(ctx)
	})
	if err != nil {
		return nil, err
	}

	return result.([]net.IP), nil

}

func (g *nmapGateway) refreshCache(ctx context.Context) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	// Equivalent to `/usr/local/bin/nmap -sP 172.16.0.0/24, with a 1 minute timeout.
	scanner, err := nmap.NewScanner(
		nmap.WithContext(ctx),
		nmap.WithTargets(g.targets),
		nmap.WithPingScan(),
		nmap.WithTimingTemplate(nmap.TimingFastest),
		nmap.WithMaxParallelism(256),
		nmap.WithDisabledDNSResolution(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create nmap scanner: %v", err)
	}

	nmapResult, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to run nmap scan: %v", err)
	}

	ipAddresses := make([]net.IP, 0, 255)
	for _, host := range nmapResult.Hosts {
		ipAddresses = append(ipAddresses, net.ParseIP(host.Addresses[0].String()))
	}

	g.cache = ipAddresses
	go func() {
		<-time.NewTimer(_cacheDuration).C
		g.cache = nil
	}()
	return ipAddresses, nil
}
