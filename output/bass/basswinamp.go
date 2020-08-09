// +build windows,386

package bass

// #cgo CFLAGS: -Iinclude
// #cgo LDFLAGS: -Llib/windows/386 -lbass_winamp.dll
// #include "BassWinamp.h"
import "C"

const (
	BASS_CTYPE_STREAM_WINAMP         = C.BASS_CTYPE_STREAM_WINAMP
	BASS_WINAMP_SYNC_BITRATE         = C.BASS_WINAMP_SYNC_BITRATE
	BASS_CONFIG_WINAMP_INPUT_TIMEOUT = C.BASS_CONFIG_WINAMP_INPUT_TIMEOUT
	BASS_WINAMP_FIND_INPUT           = C.BASS_WINAMP_FIND_INPUT
	BASS_WINAMP_FIND_RECURSIVE       = C.BASS_WINAMP_FIND_RECURSIVE
	BASS_WINAMP_FIND_COMMALIST       = C.BASS_WINAMP_FIND_COMMALIST
)

func BASS_WINAMP_FindPlugins(pluginPath string, flags uint) string {
	p := C.BASS_WINAMP_FindPlugins(C.CString(pluginPath), C.DWORD(flags))
	return C.GoString(p)
}

func BASS_WINAMP_LoadPlugin(f string) uint {
	return uint(C.BASS_WINAMP_LoadPlugin(C.CString(f)))
}

func BASS_WINAMP_UnloadPlugin(handle uint) {
	C.BASS_WINAMP_UnloadPlugin(C.DWORD(handle))
}

func BASS_WINAMP_GetName(handle uint) string {
	p := C.BASS_WINAMP_GetName(C.DWORD(handle))
	return C.GoString(p)
}

func BASS_WINAMP_GetVersion(handle uint) int {
	return int(C.BASS_WINAMP_GetVersion(C.DWORD(handle)))
}

func BASS_WINAMP_GetIsSeekable(handle uint) bool {
	return C.BASS_WINAMP_GetIsSeekable(C.DWORD(handle)) != 0
}

func BASS_WINAMP_GetUsesOutput(handle uint) bool {
	return C.BASS_WINAMP_GetUsesOutput(C.DWORD(handle)) != 0
}

func BASS_WINAMP_GetExtentions(handle uint) string {
	p := C.BASS_WINAMP_GetExtentions(C.DWORD(handle))
	return C.GoString(p)
}

func BASS_WINAMP_InfoDlg(f string, win uint) bool {
	return C.BASS_WINAMP_InfoDlg(C.CString(f), C.DWORD(win)) != 0
}

func BASS_WINAMP_ConfigPlugin(handle uint, win uint) {
	C.BASS_WINAMP_ConfigPlugin(C.DWORD(handle), C.DWORD(win))
}

func BASS_WINAMP_AboutPlugin(handle uint, win uint) {
	C.BASS_WINAMP_AboutPlugin(C.DWORD(handle), C.DWORD(win))
}

func BASS_WINAMP_StreamCreate(f string, flags uint) uint {
	return uint(C.BASS_WINAMP_StreamCreate(C.CString(f), C.DWORD(flags)))
}
