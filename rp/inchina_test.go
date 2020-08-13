package rp

import (
	"net"
	"testing"
)

func TestInChina(t *testing.T) {
	err := LoadChinaIPList()
	if err != nil {
		t.Error(err)
	}

	b := InChina("139.219.238.126")
	if !b {
		t.Error("139.219.238.126 should be in China")
	}

	b = InChina("23.99.108.233")
	if b {
		t.Error("23.99.108.233 shouldn't be in China")
	}
}

func TestIPv4InChina(t *testing.T) {
	err := LoadChinaIPList()
	if err != nil {
		t.Error(err)
	}

	ip := net.ParseIP("139.219.238.126").To4()
	b := IPv4InChina(ip)
	if !b {
		t.Error("139.219.238.126 should be in China")
	}

	ip = net.ParseIP("23.99.108.233").To4()
	b = IPv4InChina(ip)
	if b {
		t.Error("23.99.108.233 shouldn't be in China")
	}
}
