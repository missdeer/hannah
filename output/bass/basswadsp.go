// +build windows,386

package bass

// #cgo CFLAGS: -Iinclude
// #cgo LDFLAGS: -Llib/windows/386 -lBASS_WADSP.dll
// #include "BASS_WADSP.h"
import "C"
import (
	"unsafe"
)

const (
	BASS_WADSP_IPC_GETOUTPUTTIME    = C.BASS_WADSP_IPC_GETOUTPUTTIME
	BASS_WADSP_IPC_ISPLAYING        = C.BASS_WADSP_IPC_ISPLAYING
	BASS_WADSP_IPC_GETVERSION       = C.BASS_WADSP_IPC_GETVERSION
	BASS_WADSP_IPC_STARTPLAY        = C.BASS_WADSP_IPC_STARTPLAY
	BASS_WADSP_IPC_GETINFO          = C.BASS_WADSP_IPC_GETINFO
	BASS_WADSP_IPC_GETLISTLENGTH    = C.BASS_WADSP_IPC_GETLISTLENGTH
	BASS_WADSP_IPC_GETLISTPOS       = C.BASS_WADSP_IPC_GETLISTPOS
	BASS_WADSP_IPC_GETPLAYLISTFILE  = C.BASS_WADSP_IPC_GETPLAYLISTFILE
	BASS_WADSP_IPC_GETPLAYLISTTITLE = C.BASS_WADSP_IPC_GETPLAYLISTTITLE
	BASS_WADSP_IPC                  = C.BASS_WADSP_IPC
)

func BASS_WADSP_GetVersion() int {
	return int(C.BASS_WADSP_GetVersion())
}

func BASS_WADSP_Init(hwndMain uint) bool {
	return C.BASS_WADSP_Init(C.HWND(hwndMain)) != 0
}

func BASS_WADSP_Free() bool {
	return C.BASS_WADSP_Free() != 0
}

func BASS_WADSP_FreeDSP(plugin uint) bool {
	return C.BASS_WADSP_FreeDSP(C.HWADSP(plugin)) != 0
}

func BASS_WADSP_GetFakeWinampWnd(plugin uint) uint {
	return uint(C.BASS_WADSP_GetFakeWinampWnd(C.HWADSP(plugin)))
}

func BASS_WADSP_SetSongTitle(plugin uint, title string) bool {
	return C.BASS_WADSP_SetSongTitle(C.HWADSP(plugin), C.CString(title)) != 0
}

func BASS_WADSP_SetFileName(plugin uint, f string) bool {
	return C.BASS_WADSP_SetFileName(C.HWADSP(plugin), C.CString(f)) != 0
}

func BASS_WADSP_Load(f string, x int, y int, width int, height int, proc *C.WINAMPWINPROC) uint {
	return uint(C.BASS_WADSP_Load(C.CString(f), C.int(x), C.int(y), C.int(width), C.int(height), proc))
}

func BASS_WADSP_Config(plugin uint) bool {
	return C.BASS_WADSP_Config(C.HWADSP(plugin)) != 0
}

func BASS_WADSP_Start(plugin uint, module uint, hchan uint) bool {
	return C.BASS_WADSP_Start(C.HWADSP(plugin), C.DWORD(module), C.DWORD(hchan)) != 0
}

func BASS_WADSP_Stop(plugin uint) bool {
	return C.BASS_WADSP_Stop(C.HWADSP(plugin)) != 0
}

func BASS_WADSP_SetChannel(plugin uint, hchan uint) bool {
	return C.BASS_WADSP_SetChannel(C.HWADSP(plugin), C.DWORD(hchan)) != 0
}

func BASS_WADSP_GetModule(plugin uint) uint {
	return uint(C.BASS_WADSP_GetModule(C.HWADSP(plugin)))
}

func BASS_WADSP_ChannelSetDSP(plugin uint, hchan uint, priority int) uint {
	return uint(C.BASS_WADSP_ChannelSetDSP(C.HWADSP(plugin), C.DWORD(hchan), C.int(priority)))
}

func BASS_WADSP_ChannelRemoveDSP(plugin uint) bool {
	return C.BASS_WADSP_ChannelRemoveDSP(C.HWADSP(plugin)) != 0
}

func BASS_WADSP_ModifySamplesSTREAM(plugin uint, buffer []byte) uint {
	c := len(buffer)
	b := C.malloc(C.size_t(c))
	defer C.free(b)
	var p *C.char = (*C.char)(b)
	for i := 0; i < c; i++ {
		*p = C.char(buffer[i])
		p = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	}
	return uint(C.BASS_WADSP_ModifySamplesSTREAM(C.HWADSP(plugin), b, C.DWORD(c)))
}

func BASS_WADSP_ModifySamplesDSP(plugin uint, buffer []byte) uint {
	c := len(buffer)
	b := C.malloc(C.size_t(c))
	defer C.free(b)
	var p *C.char = (*C.char)(b)
	for i := 0; i < c; i++ {
		*p = C.char(buffer[i])
		p = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	}
	return uint(C.BASS_WADSP_ModifySamplesDSP(C.HWADSP(plugin), b, C.DWORD(c)))
}

func BASS_WADSP_GetName(plugin uint) string {
	p := C.BASS_WADSP_GetName(C.HWADSP(plugin))
	return C.GoString(p)
}

func BASS_WADSP_GetModuleCount(plugin uint) uint {
	return uint(C.BASS_WADSP_GetModuleCount(C.HWADSP(plugin)))
}

func BASS_WADSP_GetModuleName(plugin uint, module uint) string {
	p := C.BASS_WADSP_GetModuleName(C.HWADSP(plugin), C.DWORD(module))
	return C.GoString(p)
}

func BASS_WADSP_PluginInfoFree() bool {
	return C.BASS_WADSP_PluginInfoFree() != 0
}

func BASS_WADSP_PluginInfoLoad(dspfile string) bool {
	return C.BASS_WADSP_PluginInfoLoad(C.CString(dspfile)) != 0
}

func BASS_WADSP_PluginInfoGetName() string {
	p := C.BASS_WADSP_PluginInfoGetName()
	return C.GoString(p)
}

func BASS_WADSP_PluginInfoGetModuleCount() uint {
	return uint(C.BASS_WADSP_PluginInfoGetModuleCount())
}

func BASS_WADSP_PluginInfoGetModuleName(module uint) string {
	p := C.BASS_WADSP_PluginInfoGetModuleName(C.HWADSP(module))
	return C.GoString(p)
}
