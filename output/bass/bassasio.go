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
	if unicode {
		return C.BASS_ASIO_SetUnicode(1) != 0
	}
	return C.BASS_ASIO_SetUnicode(0) != 0
}

func BASS_ASIO_ErrorGetCode() uint {
	return uint(C.BASS_ASIO_ErrorGetCode())
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
	if lock {
		return C.BASS_ASIO_Lock(1) != 0
	}
	return C.BASS_ASIO_Lock(0) != 0
}

func BASS_ASIO_SetNotify(proc *C.ASIONOTIFYPROC, user unsafe.Pointer) bool {
	return C.BASS_ASIO_SetNotify(proc, user) != 0
}

func BASS_ASIO_ControlPanel() bool {
	return C.BASS_ASIO_ControlPanel() != 0
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
	if input {
		return uint(C.BASS_ASIO_GetLatency(1))
	}
	return uint(C.BASS_ASIO_GetLatency(0))
}

func BASS_ASIO_GetCPU() float32 {
	return float32(C.BASS_ASIO_GetCPU())
}

func BASS_ASIO_Monitor(input int, output uint, gain uint, state uint, pan uint) bool {
	return C.BASS_ASIO_Monitor(C.int(input), C.DWORD(output), C.DWORD(gain), C.DWORD(state), C.DWORD(pan)) != 0
}

func BASS_ASIO_SetDSD(dsd bool) bool {
	if dsd {
		return C.BASS_ASIO_SetDSD(1) != 0
	}
	return C.BASS_ASIO_SetDSD(0) != 0
}

func BASS_ASIO_Future(selector uint, param unsafe.Pointer) bool {
	return C.BASS_ASIO_Future(C.DWORD(selector), param) != 0
}

func BASS_ASIO_ChannelReset(input bool, channel int, flags uint) bool {
	if input {
		return C.BASS_ASIO_ChannelReset(1, C.int(channel), C.DWORD(flags)) != 0
	}
	return C.BASS_ASIO_ChannelReset(0, C.int(channel), C.DWORD(flags)) != 0
}

func BASS_ASIO_ChannelEnable(input bool, channel uint, proc *C.ASIOPROC, user unsafe.Pointer) bool {
	if input {
		return C.BASS_ASIO_ChannelEnable(1, C.DWORD(channel), proc, user) != 0
	}
	return C.BASS_ASIO_ChannelEnable(0, C.DWORD(channel), proc, user) != 0
}

func BASS_ASIO_ChannelEnableMirror(channel uint, input2 bool, channel2 uint) bool {
	if input2 {
		return C.BASS_ASIO_ChannelEnableMirror(C.DWORD(channel), 1, C.DWORD(channel2)) != 0
	}
	return C.BASS_ASIO_ChannelEnableMirror(C.DWORD(channel), 0, C.DWORD(channel2)) != 0
}

func BASS_ASIO_ChannelEnableBASS(input bool, channel uint, handle uint, join bool) bool {
	i := 0
	if input {
		i = 1
	}
	j := 0
	if join {
		j = 1
	}
	return C.BASS_ASIO_ChannelEnableBASS(C.int(i), C.DWORD(channel), C.DWORD(handle), C.int(j)) != 0
}

func BASS_ASIO_ChannelJoin(input bool, channel uint, channel2 int) bool {
	if input {
		return C.BASS_ASIO_ChannelJoin(1, C.DWORD(channel), C.int(channel2)) != 0
	}
	return C.BASS_ASIO_ChannelJoin(0, C.DWORD(channel), C.int(channel2)) != 0
}

func BASS_ASIO_ChannelPause(input bool, channel uint) bool {
	if input {
		return C.BASS_ASIO_ChannelPause(1, C.DWORD(channel)) != 0
	}
	return C.BASS_ASIO_ChannelPause(0, C.DWORD(channel)) != 0
}

func BASS_ASIO_ChannelIsActive(input bool, channel uint) uint {
	if input {
		return uint(C.BASS_ASIO_ChannelIsActive(1, C.DWORD(channel)))
	}
	return uint(C.BASS_ASIO_ChannelIsActive(0, C.DWORD(channel)))
}

func BASS_ASIO_ChannelSetFormat(input bool, channel uint, format uint) bool {
	if input {
		return C.BASS_ASIO_ChannelSetFormat(1, C.DWORD(channel), C.DWORD(format)) != 0
	}
	return C.BASS_ASIO_ChannelSetFormat(0, C.DWORD(channel), C.DWORD(format)) != 0
}

func BASS_ASIO_ChannelGetFormat(input bool, channel uint) uint {
	if input {
		return uint(C.BASS_ASIO_ChannelGetFormat(1, C.DWORD(channel)))
	}
	return uint(C.BASS_ASIO_ChannelGetFormat(0, C.DWORD(channel)))
}

func BASS_ASIO_ChannelSetRate(input bool, channel uint, rate float64) bool {
	if input {
		return C.BASS_ASIO_ChannelSetRate(1, C.DWORD(channel), C.double(rate)) != 0
	}
	return C.BASS_ASIO_ChannelSetRate(0, C.DWORD(channel), C.double(rate)) != 0
}

func BASS_ASIO_ChannelGetRate(input bool, channel uint) float64 {
	if input {
		return float64(C.BASS_ASIO_ChannelGetRate(1, C.DWORD(channel)))
	}
	return float64(C.BASS_ASIO_ChannelGetRate(0, C.DWORD(channel)))
}

func BASS_ASIO_ChannelSetVolume(input bool, channel int, volume float32) bool {
	if input {
		return C.BASS_ASIO_ChannelSetVolume(1, C.int(channel), C.float(volume)) != 0
	}
	return C.BASS_ASIO_ChannelSetVolume(0, C.int(channel), C.float(volume)) != 0
}

func BASS_ASIO_ChannelGetVolume(input bool, channel int) float32 {
	if input {
		return float32(C.BASS_ASIO_ChannelGetVolume(1, C.int(channel)))
	}
	return float32(C.BASS_ASIO_ChannelGetVolume(0, C.int(channel)))
}

func BASS_ASIO_ChannelGetLevel(input bool, channel uint) float32 {
	if input {
		return float32(C.BASS_ASIO_ChannelGetLevel(1, C.DWORD(channel)))
	}
	return float32(C.BASS_ASIO_ChannelGetLevel(0, C.DWORD(channel)))
}
