package entities

import (
	"net"
	"time"
)

const (
	DNSStatusNew   = "NEW"
	DNSStatusReady = "READY"
	DNSStatusError = "ERROR"
)

type DNS struct {
	Domain  string    `json:"domain"`
	IPV4    net.IP    `json:"ipv4"`
	Status  string    `json:"status"`
	Error   string    `json:"error"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (d *DNS) NeedsAdapting() bool {
	return d.Status == DNSStatusNew
}
