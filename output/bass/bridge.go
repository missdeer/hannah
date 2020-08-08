package bass

/*
#cgo CPPFLAGS: -Iinclude
#cgo CXXFLAGS: -Iinclude
#include "bass.h"
extern void onBASSSyncEnd(HSYNC handle, DWORD channel, DWORD data, void *user);

void cgoOnBASSSyncEnd(HSYNC handle, DWORD channel, DWORD data, void *user) {
	onBASSSyncEnd(handle, channel, data, user);
}
*/
import "C"

func bool2Cint(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}
