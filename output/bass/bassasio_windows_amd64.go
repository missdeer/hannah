// +build windows,amd64

package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #cgo LDFLAGS: -Llib/windows -lbassasio.dll
// #include "bassasio.h"
import "C"
