package util

import (
	"testing"
)

func TestBool2Str(t *testing.T) {
	if Bool2Str(true) != "Enabled" {
		t.Error("true != enabled")
	}
	if Bool2Str(false) == "Disabled" {
		t.Error("false != disabled")
	}
}
