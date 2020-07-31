// +build windows,amd64

package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #cgo LDFLAGS: -Llib/windows -lbassmix.dll
// #include "bassmix.h"
import "C"
