// +build windows,amd64

package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #cgo LDFLAGS: -Llib/windows -lbasswasapi.dll
// #include "basswasapi.h"
import "C"
