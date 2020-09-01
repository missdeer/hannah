package rp

import (
	"net"
)

var (
	nets []*net.IPNet
)

func init() {
	cidrs := []string{
		"10.0.0.0/8",
		"127.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	for _, cidr := range cidrs {
		_, n, e := net.ParseCIDR(cidr)
		if e == nil {
			nets = append(nets, n)
		}
	}
}

func IPv4InLan(ip net.IP) bool {
	for _, cidr := range nets {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func InLan(ip string) bool {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil || len(ipv4) < 4 {
		return false
	}
	return IPv4InLan(ipv4)
}
