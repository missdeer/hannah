package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #include "bass.h"
import "C"
import (
	"log"
	"unsafe"
)

const ( // Define some constant of bass
	BASS_ACTIVE_STOPPED = C.BASS_ACTIVE_STOPPED
	BASS_ACTIVE_PLAYING = C.BASS_ACTIVE_PLAYING
	BASS_ACTIVE_PAUSED  = C.BASS_ACTIVE_PAUSED
	BASS_ACTIVE_STALLED = C.BASS_ACTIVE_STALLED
)

// BASS POSsud
const (
	BASS_POS_BYTE        = C.BASS_POS_BYTE
	BASS_POS_MUSIC_ORDER = C.BASS_POS_MUSIC_ORDER
	BASS_POS_DECODE      = C.BASS_POS_DECODE
	BASS_POS_OGG         = C.BASS_POS_OGG
)

// Attribs
const (
	// Channel attributes
	BASS_ATTRIB_FREQ             = 1
	BASS_ATTRIB_VOL              = 2
	BASS_ATTRIB_PAN              = 3
	BASS_ATTRIB_EAXMIX           = 4
	BASS_ATTRIB_NOBUFFER         = 5
	BASS_ATTRIB_VBR              = 6
	BASS_ATTRIB_CPU              = 7
	BASS_ATTRIB_SRC              = 8
	BASS_ATTRIB_NET_RESUME       = 9
	BASS_ATTRIB_SCANINFO         = 10
	BASS_ATTRIB_NORAMP           = 11
	BASS_ATTRIB_BITRATE          = 12
	BASS_ATTRIB_BUFFER           = 13
	BASS_ATTRIB_MUSIC_AMPLIFY    = 0x100
	BASS_ATTRIB_MUSIC_PANSEP     = 0x101
	BASS_ATTRIB_MUSIC_PSCALER    = 0x102
	BASS_ATTRIB_MUSIC_BPM        = 0x103
	BASS_ATTRIB_MUSIC_SPEED      = 0x104
	BASS_ATTRIB_MUSIC_VOL_GLOBAL = 0x105
	BASS_ATTRIB_MUSIC_ACTIVE     = 0x106
	BASS_ATTRIB_MUSIC_VOL_CHAN   = 0x200 // + channel #
	BASS_ATTRIB_MUSIC_VOL_INST   = 0x300 // + instrument #
)

// ERRORS
const (
	// Error codes returned by BASS_ErrorGetCode
	BASS_OK             = 0  // all is OK
	BASS_ERROR_MEM      = 1  // memory error
	BASS_ERROR_FILEOPEN = 2  // can't open the file
	BASS_ERROR_DRIVER   = 3  // can't find a free/valid driver
	BASS_ERROR_BUFLOST  = 4  // the sample buffer was lost
	BASS_ERROR_HANDLE   = 5  // invalid handle
	BASS_ERROR_FORMAT   = 6  // unsupported sample format
	BASS_ERROR_POSITION = 7  // invalid position
	BASS_ERROR_INIT     = 8  // BASS_Init has not been successfully called
	BASS_ERROR_START    = 9  // BASS_Start has not been successfully called
	BASS_ERROR_SSL      = 10 // SSL/HTTPS support isn't available
	BASS_ERROR_ALREADY  = 14 // already initialized/paused/whatever
	BASS_ERROR_NOCHAN   = 18 // can't get a free channel
	BASS_ERROR_ILLTYPE  = 19 // an illegal type was specified
	BASS_ERROR_ILLPARAM = 20 // an illegal parameter was specified
	BASS_ERROR_NO3D     = 21 // no 3D support
	BASS_ERROR_NOEAX    = 22 // no EAX support
	BASS_ERROR_DEVICE   = 23 // illegal device number
	BASS_ERROR_NOPLAY   = 24 // not playing
	BASS_ERROR_FREQ     = 25 // illegal sample rate
	BASS_ERROR_NOTFILE  = 27 // the stream is not a file stream
	BASS_ERROR_NOHW     = 29 // no hardware voices available
	BASS_ERROR_EMPTY    = 31 // the MOD music has no sequence data
	BASS_ERROR_NONET    = 32 // no internet connection could be opened
	BASS_ERROR_CREATE   = 33 // couldn't create the file
	BASS_ERROR_NOFX     = 34 // effects are not available
	BASS_ERROR_NOTAVAIL = 37 // requested data/action is not available
	BASS_ERROR_DECODE   = 38 // the channel is/isn't a "decoding channel"
	BASS_ERROR_DX       = 39 // a sufficient DirectX version is not installed
	BASS_ERROR_TIMEOUT  = 40 // connection timedout
	BASS_ERROR_FILEFORM = 41 // unsupported file format
	BASS_ERROR_SPEAKER  = 42 // unavailable speaker
	BASS_ERROR_VERSION  = 43 // invalid BASS version (used by add-ons)
	BASS_ERROR_CODEC    = 44 // codec is not available/supported
	BASS_ERROR_ENDED    = 45 // the channel/file has ended
	BASS_ERROR_BUSY     = 46 // the device is busy
	BASS_ERROR_UNKNOWN  = -1 // some other mystery problem
)

// Sync type
const (
	BASS_SYNC_POS        = C.BASS_SYNC_POS
	BASS_SYNC_END        = C.BASS_SYNC_END
	BASS_SYNC_META       = C.BASS_SYNC_META
	BASS_SYNC_SLIDE      = C.BASS_SYNC_SLIDE
	BASS_SYNC_STALL      = C.BASS_SYNC_STALL
	BASS_SYNC_DOWNLOAD   = C.BASS_SYNC_DOWNLOAD
	BASS_SYNC_FREE       = C.BASS_SYNC_FREE
	BASS_SYNC_SETPOS     = C.BASS_SYNC_SETPOS
	BASS_SYNC_MUSICPOS   = C.BASS_SYNC_MUSICPOS
	BASS_SYNC_MUSICINST  = C.BASS_SYNC_MUSICINST
	BASS_SYNC_MUSICFX    = C.BASS_SYNC_MUSICFX
	BASS_SYNC_OGG_CHANGE = C.BASS_SYNC_OGG_CHANGE
	BASS_SYNC_DEV_FAIL   = C.BASS_SYNC_DEV_FAIL
	BASS_SYNC_DEV_FORMAT = C.BASS_SYNC_DEV_FORMAT
	BASS_SYNC_THREAD     = C.BASS_SYNC_THREAD
	BASS_SYNC_MIXTIME    = C.BASS_SYNC_MIXTIME
	BASS_SYNC_ONETIME    = C.BASS_SYNC_ONETIME
)

type ulong C.ulong

// ------Initialization,into,tec..
// BassInit Initializes an output device.
func Init() int {
	return int(C.BASS_Init(-1, 44100, 0, nil, nil))
}

// Free Frees all resources used by the output device, including all its samples, streams and MOD musics.
func Free() int {
	return int(C.BASS_Free())
}

// ErrorGetCode Retrieves the error code for the most recent BASS function call in the current thread.
func ErrorGetCode() int {
	return int(C.BASS_ErrorGetCode())
}

// GetVersion Retrieves the version of BASS that is loaded.
func GetVersion() int {
	return int(C.BASS_GetVersion())
}

// GetVolume Retrieves the current master volumeBase level.
func GetVolume() float32 {
	return float32(C.BASS_GetVolume())
}

// GetInfo Retrieves information on the device being used.
func GetInfo(info *C.BASS_INFO) int {
	return int(C.BASS_GetInfo(info))
}

// Pause Stops the output, pausing all musics/samples/streams on it.
func Pause() int {
	return int(C.BASS_Pause())
}

// SetDevice Sets the device to use for subsequent calls in the current thread.
func SetDevice(device C.uint) int {
	return int(C.BASS_SetDevice(C.DWORD(device)))
}

func SetVolume(value float32) bool {
	isok := int(C.BASS_SetVolume(C.float(value)))
	if isok == 0 {
		return false
	}
	return true
}

func SetChanAttr(handle uint, attr uint, value float32) uint {
	r := uint(C.BASS_ChannelSetAttribute(C.DWORD(handle), C.DWORD(attr), C.float(value)))
	return r
}

func GetChanAttr(handle uint, attr uint) float32 {
	var value float32
	C.BASS_ChannelGetAttribute(C.DWORD(handle), C.DWORD(attr), (*C.float)(unsafe.Pointer(&value)))
	return value
}

func GetChanVol(handle uint) uint {
	return uint(GetChanAttr(handle, BASS_ATTRIB_VOL) * 100)
}
func SetChanVol(handle uint, value uint) uint {
	return SetChanAttr(handle, BASS_ATTRIB_VOL, float32(value)/100)
}

// GetDevice Retrieves the device setting of the current thread.
func GetDevice(device C.ulong) int {
	return int(C.BASS_GetDevice())
}

// GetCPU Retrieves the current CPU usage of BASS.
func GetCPU() float32 {
	return float32(C.BASS_GetCPU())
}

// --------------------------------------

// ------Streams-------------------------

// StreamCreate Creates a user sample stream.
func StreamCreate(freq uint, proc *C.STREAMPROC, user unsafe.Pointer) uint {
	return uint(C.BASS_StreamCreate(C.DWORD(freq), 2, C.BASS_SAMPLE_FLOAT, proc, user))
}

// StreamCreateFile Creates a sample stream from an MP3, MP2, MP1, OGG, WAV, AIFF or plugin supported file.
func StreamCreateFile(mem int, file string, offset uint64, length uint64) uint {
	return uint(C.BASS_StreamCreateFile(C.int(mem), unsafe.Pointer(C.CString(file)), C.ulonglong(offset), C.ulonglong(length), C.DWORD(C.BASS_SAMPLE_FLOAT)))
}

// StreamCreateURL ates a sample stream from an MP3, MP2, MP1, OGG, WAV, AIFF or plugin supported file on the internet, optionally receiving the downloaded data in a callback function.
func StreamCreateURL(url string, offset uint, proc *C.DOWNLOADPROC, user unsafe.Pointer) uint {
	return uint(C.BASS_StreamCreateURL(C.CString(url), C.DWORD(offset), C.BASS_SAMPLE_FLOAT, proc, user))
}

func StreamPutData(handle uint, buffer []byte, length int) uint32 {

	return uint32(C.BASS_StreamPutData(C.DWORD(handle), C.CBytes(buffer), C.DWORD(length)))
}

func StreamFree(handle uint) uint32 {
	return uint32(C.BASS_StreamFree(C.DWORD(handle)))
}

// --------------------------------------

// ------Channels-------------------------

// ChannelPlay Starts (or resumes) playback of a sample, stream, MOD music, or recording.
func ChannelPlay(handle uint, restart int) int {
	return int(C.BASS_ChannelPlay(C.DWORD(handle), C.int(restart)))
}

// ChannelPause a sample, stream, MOD music, or recording.
func ChannelPause(handle uint) int {
	return int(C.BASS_ChannelPause(C.DWORD(handle)))
}

// ChannelStop Stops a sample, stream, MOD music, or recording.
func ChannelStop(handle uint) int {
	return int(C.BASS_ChannelStop(C.DWORD(handle)))
}

// ChannelBytes2Seconds Translates a byte position into time (seconds), based on a channel's format.
func ChannelBytes2Seconds(handle uint, pos int) int {
	return int(C.BASS_ChannelBytes2Seconds(C.DWORD(handle), C.ulonglong(pos)))
}

// ChannelSeconds2Bytes Translates a time (seconds) position into bytes, based on a channel's format.
func ChannelSeconds2Bytes(handle uint, pos int) int {
	return int(C.BASS_ChannelSeconds2Bytes(C.DWORD(handle), C.double(pos)))
}

// ChannelIsActive Checks if a sample, stream, or MOD music is active (playing) or stalled. Can also check if a recording is in progress.
func ChannelIsActive(handle uint) int {
	return int(C.BASS_ChannelIsActive(C.DWORD(handle)))
}

// ChannelGetPosition Retrieves the playback position of a sample, stream, or MOD music. Can also be used with a recording channel.
func ChannelGetPosition(handle uint, mode int) int {
	return int(C.BASS_ChannelGetPosition(C.DWORD(handle), C.DWORD(mode)))
}

// ChannelSetPosition Sets the playback position of a sample, MOD music, or stream.
func ChannelSetPosition(handle uint, pos int, mode int) int {
	return int(C.BASS_ChannelSetPosition(C.DWORD(handle), C.ulonglong(pos), C.DWORD(mode)))
}

// ChannelSetAttribute Sets the value of a channel's attribute.
func ChannelSetAttribute(handle uint, attrib C.uint, value C.float) int {
	return int(C.BASS_ChannelSetAttribute(C.DWORD(handle), C.DWORD(attrib), value))
}

// ChannelUpdate Updates the playback buffer of a stream or MOD music.
func ChannelUpdate(handle uint, length C.uint) int {
	return int(C.BASS_ChannelUpdate(C.DWORD(handle), C.DWORD(length)))
}

// ChannelGetLength Retrieves the playback length of a channel.
func ChannelGetLength(handle uint, mode int) int {
	return int(C.BASS_ChannelGetLength(C.DWORD(handle), C.DWORD(mode)))
}

// ChannelSetSync set sync callback for a channel
func ChannelSetSync(handle uint, syncType uint, param uint64, proc *C.SYNCPROC, user unsafe.Pointer) uint {
	return uint(C.BASS_ChannelSetSync(C.DWORD(handle), C.DWORD(syncType), C.QWORD(param), proc, user))
}

// ---------------------------------------
const (
	BASS_UNICODE = C.BASS_UNICODE
)

// PluginLoad load a bass plugin
func PluginLoad(file string) int {
	result := int(C.BASS_PluginLoad(C.CString(file), C.DWORD(0)))
	switch result {
	case BASS_ERROR_FILEOPEN:
		log.Println(file, "err: BASS_ERROR_FILEOPEN")
		break
	case BASS_ERROR_FILEFORM:
		log.Println(file, "err: BASS_ERROR_FILEFORM")
		break
	case BASS_ERROR_VERSION:
		log.Println(file, "err: BASS_ERROR_VERSION")
		break
	}
	return result
}

// PluginFree free a bass plugin
func PluginFree(handle int) bool {
	result := int(C.BASS_PluginFree(C.DWORD(handle)))
	if result == 0 {
		return false
	}
	return true
}
