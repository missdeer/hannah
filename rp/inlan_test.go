package rp

import (
	"testing"
)

func TestInLan(t *testing.T) {
	if !InLan("127.0.0.1") {
		t.Error("127.0.0.1")
	}
	if !InLan("192.168.1.1") {
		t.Error("192.168.1.1")
	}
	if !InLan("10.140.12.1") {
		t.Error("10.140.12.1")
	}

	if InLan("192.5.6.30") {
		t.Error("192.5.6.30")
	}
}
