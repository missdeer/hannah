// +build windows

package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #include "bassasio.h"
import "C"
import (
	"unsafe"
)

const (
	BASSASIOVERSION = C.BASSASIOVERSION
)

// BASS_ASIO_Init flags
const (
	BASS_ASIO_THREAD    = C.BASS_ASIO_THREAD
	BASS_ASIO_JOINORDER = C.BASS_ASIO_JOINORDER
)

// sample formats
const (
	BASS_ASIO_FORMAT_16BIT   = C.BASS_ASIO_FORMAT_16BIT
	BASS_ASIO_FORMAT_24BIT   = C.BASS_ASIO_FORMAT_24BIT
	BASS_ASIO_FORMAT_32BIT   = C.BASS_ASIO_FORMAT_32BIT
	BASS_ASIO_FORMAT_FLOAT   = C.BASS_ASIO_FORMAT_FLOAT
	BASS_ASIO_FORMAT_DSD_LSB = C.BASS_ASIO_FORMAT_DSD_LSB
	BASS_ASIO_FORMAT_DSD_MSB = C.BASS_ASIO_FORMAT_DSD_MSB
	BASS_ASIO_FORMAT_DITHER  = C.BASS_ASIO_FORMAT_DITHER
)

// BASS_ASIO_ChannelReset flags
const (
	BASS_ASIO_RESET_ENABLE = C.BASS_ASIO_RESET_ENABLE
	BASS_ASIO_RESET_JOIN   = C.BASS_ASIO_RESET_JOIN
	BASS_ASIO_RESET_PAUSE  = C.BASS_ASIO_RESET_PAUSE
	BASS_ASIO_RESET_FORMAT = C.BASS_ASIO_RESET_FORMAT
	BASS_ASIO_RESET_RATE   = C.BASS_ASIO_RESET_RATE
	BASS_ASIO_RESET_VOLUME = C.BASS_ASIO_RESET_VOLUME
	BASS_ASIO_RESET_JOINED = C.BASS_ASIO_RESET_JOINED
)

// BASS_ASIO_ChannelIsActive return values
const (
	BASS_ASIO_ACTIVE_DISABLED = C.BASS_ASIO_ACTIVE_DISABLED
	BASS_ASIO_ACTIVE_ENABLED  = C.BASS_ASIO_ACTIVE_ENABLED
	BASS_ASIO_ACTIVE_PAUSED   = C.BASS_ASIO_ACTIVE_PAUSED
)

// driver notifications
const (
	BASS_ASIO_NOTIFY_RATE  = C.BASS_ASIO_NOTIFY_RATE
	BASS_ASIO_NOTIFY_RESET = C.BASS_ASIO_NOTIFY_RESET
)

// BASS_ASIO_ChannelGetLevel flags
const (
	BASS_ASIO_LEVEL_RMS = C.BASS_ASIO_LEVEL_RMS
)

func BASS_ASIO_GetVersion() uint {
	return uint(C.BASS_ASIO_GetVersion())
}

func BASS_ASIO_SetUnicode(unicode bool) bool {
	return C.BASS_ASIO_SetUnicode(bool2Cint(unicode)) != 0
}

func BASS_ASIO_ErrorGetCode() uint {
	return uint(C.BASS_ASIO_ErrorGetCode())
}

type BassASIODeviceInfo struct {
	Name   string
	Driver string
}

func BASS_ASIO_GetDeviceInfo(device uint, info *BassASIODeviceInfo) bool {
	i := (*C.BASS_ASIO_DEVICEINFO)(C.malloc(C.sizeof_BASS_ASIO_DEVICEINFO))
	defer C.free(unsafe.Pointer(i))
	res := C.BASS_ASIO_GetDeviceInfo(C.DWORD(device), i) != 0
	if res {
		info.Driver = C.GoString(i.driver)
		info.Name = C.GoString(i.name)
	}
	return res
}

func BASS_ASIO_SetDevice(device uint) bool {
	return C.BASS_ASIO_SetDevice(C.DWORD(device)) != 0
}

func BASS_ASIO_GetDevice() uint {
	return uint(C.BASS_ASIO_GetDevice())
}

func BASS_ASIO_Init(device int, flags uint) bool {
	return C.BASS_ASIO_Init(C.int(device), C.DWORD(flags)) != 0
}

func BASS_ASIO_Free() bool {
	return C.BASS_ASIO_Free() != 0
}

func BASS_ASIO_Lock(lock bool) bool {
	return C.BASS_ASIO_Lock(bool2Cint(lock)) != 0
}

func BASS_ASIO_SetNotify(proc *C.ASIONOTIFYPROC, user unsafe.Pointer) bool {
	return C.BASS_ASIO_SetNotify(proc, user) != 0
}

func BASS_ASIO_ControlPanel() bool {
	return C.BASS_ASIO_ControlPanel() != 0
}

type BassASIOInfo struct {
	Name      string
	Version   uint
	Inputs    uint
	Outputs   uint
	BufMin    uint
	BufMax    uint
	BufPref   uint
	BufGran   int
	InitFlags uint
}

func BASS_ASIO_GetInfo(info *BassASIOInfo) bool {
	i := (*C.BASS_ASIO_INFO)(C.malloc(C.sizeof_BASS_ASIO_INFO))
	defer C.free(unsafe.Pointer(i))
	res := C.BASS_ASIO_GetInfo(i) != 0
	if res {
		info.InitFlags = uint(i.initflags)
		info.Version = uint(i.version)
		info.Inputs = uint(i.inputs)
		info.Outputs = uint(i.outputs)
		info.BufMin = uint(i.bufmin)
		info.BufMax = uint(i.bufmax)
		info.BufPref = uint(i.bufpref)
		info.InitFlags = uint(i.initflags)
		info.BufGran = int(i.bufgran)
		info.Name = C.GoStringN(&i.name[0], 32)
	}
	return res
}

func BASS_ASIO_CheckRate(rate float64) bool {
	return C.BASS_ASIO_CheckRate(C.double(rate)) != 0
}

func BASS_ASIO_SetRate(rate float64) bool {
	return C.BASS_ASIO_SetRate(C.double(rate)) != 0
}

func BASS_ASIO_GetRate() float64 {
	return float64(C.BASS_ASIO_GetRate())
}

func BASS_ASIO_Start(buflen uint, threads uint) bool {
	return C.BASS_ASIO_Start(C.DWORD(buflen), C.DWORD(threads)) != 0
}

func BASS_ASIO_Stop() bool {
	return C.BASS_ASIO_Stop() != 0
}

func BASS_ASIO_IsStarted() bool {
	return C.BASS_ASIO_IsStarted() != 0
}

func BASS_ASIO_GetLatency(input bool) uint {
	return uint(C.BASS_ASIO_GetLatency(bool2Cint(input)))
}

func BASS_ASIO_GetCPU() float32 {
	return float32(C.BASS_ASIO_GetCPU())
}

func BASS_ASIO_Monitor(input int, output uint, gain uint, state uint, pan uint) bool {
	return C.BASS_ASIO_Monitor(C.int(input), C.DWORD(output), C.DWORD(gain), C.DWORD(state), C.DWORD(pan)) != 0
}

func BASS_ASIO_SetDSD(dsd bool) bool {
	return C.BASS_ASIO_SetDSD(bool2Cint(dsd)) != 0
}

func BASS_ASIO_Future(selector uint, param unsafe.Pointer) bool {
	return C.BASS_ASIO_Future(C.DWORD(selector), param) != 0
}

type BassASIOChannelInfo struct {
	Group  uint
	Format uint
	Name   string
}

func BASS_ASIO_ChannelGetInfo(input bool, channel uint, info *BassASIOChannelInfo) bool {
	i := (*C.BASS_ASIO_CHANNELINFO)(C.malloc(C.sizeof_BASS_ASIO_CHANNELINFO))
	defer C.free(unsafe.Pointer(i))
	res := C.BASS_ASIO_ChannelGetInfo(bool2Cint(input), C.DWORD(channel), i) != 0
	if res {
		info.Group = uint(i.group)
		info.Format = uint(i.format)
		info.Name = C.GoStringN(&i.name[0], 32)
	}
	return res
}

func BASS_ASIO_ChannelReset(input bool, channel int, flags uint) bool {
	return C.BASS_ASIO_ChannelReset(bool2Cint(input), C.int(channel), C.DWORD(flags)) != 0
}

func BASS_ASIO_ChannelEnable(input bool, channel uint, proc *C.ASIOPROC, user unsafe.Pointer) bool {
	return C.BASS_ASIO_ChannelEnable(bool2Cint(input), C.DWORD(channel), proc, user) != 0
}

func BASS_ASIO_ChannelEnableMirror(channel uint, input2 bool, channel2 uint) bool {
	return C.BASS_ASIO_ChannelEnableMirror(C.DWORD(channel), bool2Cint(input2), C.DWORD(channel2)) != 0
}

func BASS_ASIO_ChannelEnableBASS(input bool, channel uint, handle uint, join bool) bool {
	return C.BASS_ASIO_ChannelEnableBASS(bool2Cint(input), C.DWORD(channel), C.DWORD(handle), bool2Cint(join)) != 0
}

func BASS_ASIO_ChannelJoin(input bool, channel uint, channel2 int) bool {
	return C.BASS_ASIO_ChannelJoin(bool2Cint(input), C.DWORD(channel), C.int(channel2)) != 0
}

func BASS_ASIO_ChannelPause(input bool, channel uint) bool {
	return C.BASS_ASIO_ChannelPause(bool2Cint(input), C.DWORD(channel)) != 0
}

func BASS_ASIO_ChannelIsActive(input bool, channel uint) uint {
	return uint(C.BASS_ASIO_ChannelIsActive(bool2Cint(input), C.DWORD(channel)))
}

func BASS_ASIO_ChannelSetFormat(input bool, channel uint, format uint) bool {
	return C.BASS_ASIO_ChannelSetFormat(bool2Cint(input), C.DWORD(channel), C.DWORD(format)) != 0
}

func BASS_ASIO_ChannelGetFormat(input bool, channel uint) uint {
	return uint(C.BASS_ASIO_ChannelGetFormat(bool2Cint(input), C.DWORD(channel)))
}

func BASS_ASIO_ChannelSetRate(input bool, channel uint, rate float64) bool {
	return C.BASS_ASIO_ChannelSetRate(bool2Cint(input), C.DWORD(channel), C.double(rate)) != 0
}

func BASS_ASIO_ChannelGetRate(input bool, channel uint) float64 {
	return float64(C.BASS_ASIO_ChannelGetRate(bool2Cint(input), C.DWORD(channel)))
}

func BASS_ASIO_ChannelSetVolume(input bool, channel int, volume float32) bool {
	return C.BASS_ASIO_ChannelSetVolume(bool2Cint(input), C.int(channel), C.float(volume)) != 0
}

func BASS_ASIO_ChannelGetVolume(input bool, channel int) float32 {
	return float32(C.BASS_ASIO_ChannelGetVolume(bool2Cint(input), C.int(channel)))
}

func BASS_ASIO_ChannelGetLevel(input bool, channel uint) float32 {
	return float32(C.BASS_ASIO_ChannelGetLevel(bool2Cint(input), C.DWORD(channel)))
}
