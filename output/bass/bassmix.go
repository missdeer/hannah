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

func BASS_Mixer_ChannelGetMatrix(handle uint) (m []float32) {
	matrix := C.malloc(8 * C.sizeof_float)
	defer C.free(matrix)
	res := C.BASS_Mixer_ChannelGetMatrix(C.DWORD(handle), matrix)
	if res == 0 {
		return nil
	}
	var p *C.float = (*C.float)(matrix)
	for i := 0; i < 8; i++ {
		m = append(m, float32(*p))
		p = (*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	}
	return m
}

type BassMixerNode struct {
	Pos uint64
	Val float32
}

func BASS_Mixer_ChannelSetEnvelope(handle uint, typ uint, nodes []BassMixerNode, count uint) bool {
	n := C.malloc(C.size_t(count * C.sizeof_BASS_MIXER_NODE))
	defer C.free(n)
	var p *C.BASS_MIXER_NODE = (*C.BASS_MIXER_NODE)(n)
	for i := 0; i < 8; i++ {
		(*p).pos = C.QWORD(nodes[i].Pos)
		(*p).value = C.float(nodes[i].Val)
		p = (*C.BASS_MIXER_NODE)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	}
	return C.BASS_Mixer_ChannelSetEnvelope(C.DWORD(handle), C.DWORD(typ), (*C.BASS_MIXER_NODE)(n), C.DWORD(count)) != 0
}

func BASS_Mixer_ChannelSetEnvelopePos(handle uint, typ uint, pos uint64) bool {
	return C.BASS_Mixer_ChannelSetEnvelopePos(C.DWORD(handle), C.DWORD(typ), C.QWORD(pos)) != 0
}

func BASS_Mixer_ChannelGetEnvelopePos(handle uint, typ uint) (value float32, pos int64) {
	v := C.malloc(C.sizeof_float)
	defer C.free(v)
	pos = int64(C.BASS_Mixer_ChannelGetEnvelopePos(C.DWORD(handle), C.DWORD(typ), (*C.float)(v)))
	if pos != -1 {
		value = float32(*(*C.float)(v))
	}
	return value, pos
}

func BASS_Split_StreamCreate(handle uint, flags uint, chanmap []int) uint {
	c := len(chanmap)
	cm := C.malloc(C.size_t(c * C.sizeof_int))
	defer C.free(cm)
	var p *C.int = (*C.int)(cm)
	for i := 0; i < c; i++ {
		*p = C.int(chanmap[i])
		p = (*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	}
	return uint(C.BASS_Split_StreamCreate(C.DWORD(handle), C.DWORD(flags), (*C.int)(cm)))
}

func BASS_Split_StreamGetSource(handle uint) uint {
	return uint(C.BASS_Split_StreamGetSource(C.HSTREAM(handle)))
}

func BASS_Split_StreamGetSplits(handle uint) (splits []uint) {
	c := C.BASS_Split_StreamGetSplits(C.DWORD(handle), nil, 0)
	s := C.malloc(C.size_t(c * C.sizeof_HSTREAM))
	defer C.free(s)
	res := C.BASS_Split_StreamGetSplits(C.DWORD(handle), (*C.HSTREAM)(s), C.DWORD(c))
	if res > 0 {
		var p *C.HSTREAM = (*C.HSTREAM)(s)
		for i := 0; i < int(res); i++ {
			splits = append(splits, uint(*p))
			p = (*C.HSTREAM)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
		}
	}
	return splits
}

func BASS_Split_StreamReset(handle uint) bool {
	return C.BASS_Split_StreamReset(C.DWORD(handle)) != 0
}

func BASS_Split_StreamResetEx(handle uint, offset uint) bool {
	return C.BASS_Split_StreamResetEx(C.DWORD(handle), C.DWORD(offset)) != 0
}

func BASS_Split_StreamGetAvailable(handle uint) uint {
	return uint(C.BASS_Split_StreamGetAvailable(C.DWORD(handle)))
}
