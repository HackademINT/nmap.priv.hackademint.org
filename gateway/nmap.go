package gateway

import (
	"context"
	"net"
	"time"
	"fmt"
	"github.com/Ullaakut/nmap"
	"golang.org/x/sync/singleflight"
)

type NmapGateway interface {
	ScanSubnet(ctx context.Context) ([]net.IP, error)
}

type nmapGateway struct {
	targets string
	group   singleflight.Group
}

func NewNmapGateway(targets string) NmapGateway {
	return &nmapGateway{
		targets: targets,
		group:   singleflight.Group{},
	}
}

func (g *nmapGateway) ScanSubnet(ctx context.Context) ([]net.IP, error) {
	result, err, _ := g.group.Do("scanSubnet", func() (interface{}, error) {
		return g.scanSubnetNotCached(ctx)
	})

	return result.([]net.IP), err

}

func (g *nmapGateway) scanSubnetNotCached(ctx context.Context) ([]net.IP, error) {
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

	return ipAddresses, nil
}
