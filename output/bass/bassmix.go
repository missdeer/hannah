// +build windows

package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #include "bassmix.h"
import "C"
import (
	"unsafe"
)

// Additional BASS_SetConfig options
const (
	BASS_CONFIG_MIXER_BUFFER = C.BASS_CONFIG_MIXER_BUFFER
	BASS_CONFIG_MIXER_POSEX  = C.BASS_CONFIG_MIXER_POSEX
	BASS_CONFIG_SPLIT_BUFFER = C.BASS_CONFIG_SPLIT_BUFFER
)

// BASS_Mixer_StreamCreate flags
const (
	BASS_MIXER_END     = C.BASS_MIXER_END
	BASS_MIXER_NONSTOP = C.BASS_MIXER_NONSTOP
	BASS_MIXER_RESUME  = C.BASS_MIXER_RESUME
	BASS_MIXER_POSEX   = C.BASS_MIXER_POSEX
)

// BASS_Mixer_StreamAddChannel/Ex flags
const (
	BASS_MIXER_CHAN_ABSOLUTE = C.BASS_MIXER_CHAN_ABSOLUTE
	BASS_MIXER_CHAN_BUFFER   = C.BASS_MIXER_CHAN_BUFFER
	BASS_MIXER_CHAN_LIMIT    = C.BASS_MIXER_CHAN_LIMIT
	BASS_MIXER_CHAN_MATRIX   = C.BASS_MIXER_CHAN_MATRIX
	BASS_MIXER_CHAN_PAUSE    = C.BASS_MIXER_CHAN_PAUSE
	BASS_MIXER_CHAN_DOWNMIX  = C.BASS_MIXER_CHAN_DOWNMIX
	BASS_MIXER_CHAN_NORAMPIN = C.BASS_MIXER_CHAN_NORAMPIN
	BASS_MIXER_BUFFER        = C.BASS_MIXER_BUFFER
	BASS_MIXER_LIMIT         = C.BASS_MIXER_LIMIT
	BASS_MIXER_MATRIX        = C.BASS_MIXER_MATRIX
	BASS_MIXER_PAUSE         = C.BASS_MIXER_PAUSE
	BASS_MIXER_DOWNMIX       = C.BASS_MIXER_DOWNMIX
	BASS_MIXER_NORAMPIN      = C.BASS_MIXER_NORAMPIN
)

// Mixer attributes
const (
	BASS_ATTRIB_MIXER_LATENCY = C.BASS_ATTRIB_MIXER_LATENCY
)

// BASS_Split_StreamCreate flags
const (
	BASS_SPLIT_SLAVE = C.BASS_SPLIT_SLAVE
	BASS_SPLIT_POS   = C.BASS_SPLIT_POS
)

// Splitter attributes
const (
	BASS_ATTRIB_SPLIT_ASYNCBUFFER = C.BASS_ATTRIB_SPLIT_ASYNCBUFFER
	BASS_ATTRIB_SPLIT_ASYNCPERIOD = C.BASS_ATTRIB_SPLIT_ASYNCPERIOD
)

// Envelope types
const (
	BASS_MIXER_ENV_FREQ   = C.BASS_MIXER_ENV_FREQ
	BASS_MIXER_ENV_VOL    = C.BASS_MIXER_ENV_VOL
	BASS_MIXER_ENV_PAN    = C.BASS_MIXER_ENV_PAN
	BASS_MIXER_ENV_LOOP   = C.BASS_MIXER_ENV_LOOP
	BASS_MIXER_ENV_REMOVE = C.BASS_MIXER_ENV_REMOVE
)

// Additional sync types
const (
	BASS_SYNC_MIXER_ENVELOPE      = C.BASS_SYNC_MIXER_ENVELOPE
	BASS_SYNC_MIXER_ENVELOPE_NODE = C.BASS_SYNC_MIXER_ENVELOPE_NODE
)

// Additional BASS_Mixer_ChannelSetPosition flag
const (
	BASS_POS_MIXER_RESET = C.BASS_POS_MIXER_RESET
)

// BASS_CHANNELINFO types
const (
	BASS_CTYPE_STREAM_MIXER = C.BASS_CTYPE_STREAM_MIXER
	BASS_CTYPE_STREAM_SPLIT = C.BASS_CTYPE_STREAM_SPLIT
)

func BASS_Mixer_GetVersion() uint {
	return uint(C.BASS_Mixer_GetVersion())
}

func BASS_Mixer_StreamCreate(freq uint, chans uint, flags uint) uint {
	return uint(C.BASS_Mixer_StreamCreate(C.DWORD(freq), C.DWORD(chans), C.DWORD(flags)))
}

func BASS_Mixer_StreamAddChannel(handle uint, channel uint, flags uint) bool {
	return C.BASS_Mixer_StreamAddChannel(C.HSTREAM(handle), C.DWORD(channel), C.DWORD(flags)) != 0
}

func BASS_Mixer_StreamAddChannelEx(handle uint, channel uint, flags uint, start uint64, length uint64) bool {
	return C.BASS_Mixer_StreamAddChannelEx(C.HSTREAM(handle), C.DWORD(channel), C.DWORD(flags), C.QWORD(start), C.QWORD(length)) != 0
}

func BASS_Mixer_StreamGetChannels(handle uint, channels []uint, count uint) int {
	c := int(C.BASS_Mixer_StreamGetChannels(C.HSTREAM(handle), nil, 0))
	chans := (*C.DWORD)(C.malloc(C.size_t(c * C.sizeof_DWORD)))
	defer C.free(unsafe.Pointer(chans))
	res := int(C.BASS_Mixer_StreamGetChannels(C.HSTREAM(handle), chans, C.DWORD(count)))
	var p *C.DWORD = chans
	for i := 0; i < int(count) && i < c; i++ {
		channels = append(channels, uint(*p))
		p = (*C.DWORD)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	}
	return res
}

func BASS_Mixer_ChannelGetMixer(handle uint) uint {
	return uint(C.BASS_Mixer_ChannelGetMixer(C.DWORD(handle)))
}

func BASS_Mixer_ChannelFlags(handle uint, flags uint, mask uint) uint {
	return uint(C.BASS_Mixer_ChannelFlags(C.DWORD(handle), C.DWORD(flags), C.DWORD(mask)))
}

func BASS_Mixer_ChannelRemove(handle uint) bool {
	return C.BASS_Mixer_ChannelRemove(C.DWORD(handle)) != 0
}

func BASS_Mixer_ChannelSetPosition(handle uint, pos uint64, mode uint) bool {
	return C.BASS_Mixer_ChannelSetPosition(C.DWORD(handle), C.QWORD(pos), C.DWORD(mode)) != 0
}

func BASS_Mixer_ChannelGetPosition(handle uint, mode uint) uint64 {
	return uint64(C.BASS_Mixer_ChannelGetPosition(C.DWORD(handle), C.DWORD(mode)))
}

func BASS_Mixer_ChannelGetPositionEx(handle uint, mode uint, delay uint) uint64 {
	return uint64(C.BASS_Mixer_ChannelGetPositionEx(C.DWORD(handle), C.DWORD(mode), C.DWORD(delay)))
}

func BASS_Mixer_ChannelGetLevel(handle uint) uint {
	return uint(C.BASS_Mixer_ChannelGetLevel(C.DWORD(handle)))
}

func BASS_Mixer_ChannelGetData(handle uint, buffer []byte, length uint) uint {
	c := int(C.BASS_Mixer_ChannelGetData(C.DWORD(handle), nil, C.BASS_DATA_AVAILABLE))
	b := C.malloc(C.size_t(c))
	defer C.free(b)
	res := uint(C.BASS_Mixer_ChannelGetData(C.DWORD(handle), b, C.DWORD(c)))
	var p *C.char = (*C.char)(b)
	for i := 0; i < c; i++ {
		buffer = append(buffer, byte(*p))
		p = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	}
	return res
}

func BASS_Mixer_ChannelSetSync(handle uint, typ uint, param uint64, proc *C.SYNCPROC, user unsafe.Pointer) uint {
	return uint(C.BASS_Mixer_ChannelSetSync(C.DWORD(handle), C.DWORD(typ), C.QWORD(param), proc, user))
}

func BASS_Mixer_ChannelRemoveSync(handle uint, sync uint) bool {
	return C.BASS_Mixer_ChannelRemoveSync(C.DWORD(handle), C.HSYNC(sync)) != 0
}

func BASS_Mixer_ChannelSetMatrix(handle uint, matrix unsafe.Pointer) bool {
	return C.BASS_Mixer_ChannelSetMatrix(C.DWORD(handle), matrix) != 0
}

func BASS_Mixer_ChannelSetMatrixEx(handle uint, matrix unsafe.Pointer, time float32) bool {
	return C.BASS_Mixer_ChannelSetMatrixEx(C.DWORD(handle), matrix, C.float(time)) != 0
}
