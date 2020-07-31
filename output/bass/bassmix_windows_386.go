// +build windows,386

package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #cgo LDFLAGS: -Llib/windows/x86 -lbassmix.dll
// #include "bassmix.h"
import "C"
