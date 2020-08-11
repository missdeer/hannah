// Package mpg123 provides an easy to use wrapper around mpg123 library
// for decoding mp3 to audio samples
package mpg123

/*
#cgo darwin windows pkg-config: libmpg123
#include <mpg123.h>
*/
import "C"

import (
	"bytes"
	"fmt"
	"log"
	"unsafe"
)

var _ = log.Print

func init() {
	C.mpg123_init()
}

// func Exit() {
// 	C.mpg123_exit()
// }

type Error int

func (e Error) Error() string {
	return C.GoString(C.mpg123_plain_strerror(C.int(e)))
}

func mpgError(err C.int) error {
	if err == C.MPG123_OK {
		return nil
	} else {
		return Error(int(err))
	}
}

var ErrDone = Error(C.MPG123_DONE)
var ErrNewFormat = Error(C.MPG123_NEW_FORMAT)
var ErrNeedMore = Error(C.MPG123_NEED_MORE)

func charPointerArray(a **C.char) []string {
	offset := unsafe.Sizeof(a)
	res := make([]string, 0)
	var i uintptr = 0
	for {
		pp := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(a)) + offset*i))
		if pp == nil {
			break
		}
		ps := C.GoString(pp)
		res = append(res, ps)
		i++

	}
	return res
}

func AvailableDecoders() []string {
	r := charPointerArray(C.mpg123_decoders())
	return r
}

type Handle struct {
	h *C.mpg123_handle
}

func newHandle(decoder *string) (*Handle, error) {
	var err C.int
	var d *C.char = nil
	if decoder != nil {
		d = C.CString(*decoder)
		defer C.free(unsafe.Pointer(d))
	}
	handle := C.mpg123_new(d, &err)
	if err != 0 {
		return nil, mpgError(err)
	}

	return &Handle{
		h: handle,
	}, nil
}

func NewHandle(decoder string) (*Handle, error) {
	return newHandle(&decoder)
}

func NewDefaultHandle() (*Handle, error) {
	return newHandle(nil)
}

func (h *Handle) Delete() {
	C.mpg123_delete(h.h)
}

func (h *Handle) CurrentDecoder() string {
	return C.GoString(C.mpg123_current_decoder(h.h))
}

func (h *Handle) OpenFeed() error {
	return mpgError(C.mpg123_open_feed(h.h))
}

func (h *Handle) Feed(data []byte) error {
	dataP := unsafe.Pointer(&data[0])
	dataLen := len(data)
	return mpgError(C.mpg123_feed(h.h, (*C.uchar)(dataP), C.size_t(dataLen)))
}

func (h *Handle) Decode(in, out []byte) (int, error) {
	var inBufC *C.uchar = nil
	var inBufCLen C.size_t = 0

	if len(in) > 0 {
		inBufC = (*C.uchar)(unsafe.Pointer(&(in[0])))
		inBufCLen = C.size_t(len(in))
	}

	var outBufC unsafe.Pointer = nil
	var outBufCLen C.size_t = 0
	if len(out) > 0 {
		outBufC = unsafe.Pointer(&(out[0]))
		outBufCLen = C.size_t(len(out))
	}
	var outBufCDone C.size_t = 0
	err := C.mpg123_decode(h.h, inBufC, inBufCLen, outBufC, outBufCLen, &outBufCDone)
	// log.Printf("Decode() %d %v", int(outBufCDone), Error(err))

	return int(outBufCDone), mpgError(err)

}

func (h *Handle) FeedSeek(sampleOffset int64, whence int) (newSampleOffset int64, inputOffset int64, err error) {
	var byteOffsetC C.off_t
	ret := C.mpg123_feedseek(h.h, C.off_t(sampleOffset), C.int(whence), &byteOffsetC)
	if ret >= 0 {
		return int64(ret), int64(byteOffsetC), nil
	} else {
		return 0, 0, mpgError(C.int(ret))
	}

}

//go:generate stringer -type=Encoding

type Encoding uint16

func (e Encoding) Size() int {
	return int(C.mpg123_encsize(C.int(e)))

}

func (e Encoding) String() string {
	s := "u"
	if uint16(e)&uint16(C.MPG123_ENC_SIGNED) != 0 {
		s = "s"
	}
	return fmt.Sprintf("%s%2d", s, e.Size()*8)
}

const (
	EncodingInt16   Encoding = Encoding(C.MPG123_ENC_SIGNED_16)
	EncodingUint16  Encoding = Encoding(C.MPG123_ENC_UNSIGNED_16)
	EncodingUint8   Encoding = Encoding(C.MPG123_ENC_UNSIGNED_8)
	EncodingInt8    Encoding = Encoding(C.MPG123_ENC_SIGNED_8)
	EncodingULaw8   Encoding = Encoding(C.MPG123_ENC_ULAW_8)
	EncodingALaw8   Encoding = Encoding(C.MPG123_ENC_ALAW_8)
	EncodingInt32   Encoding = Encoding(C.MPG123_ENC_SIGNED_32)
	EncodingUint32  Encoding = Encoding(C.MPG123_ENC_UNSIGNED_32)
	EncodingInt24   Encoding = Encoding(C.MPG123_ENC_SIGNED_24)
	EncodingUint24  Encoding = Encoding(C.MPG123_ENC_UNSIGNED_24)
	EncodingFloat32 Encoding = Encoding(C.MPG123_ENC_FLOAT_32)
	EncodingFloat64 Encoding = Encoding(C.MPG123_ENC_FLOAT_64)
	EncodingAny     Encoding = Encoding(C.MPG123_ENC_ANY)
)

// Format of the decoded mp3 audio
type OutputFormat struct {
	Rate     int
	Channels int
	Encoding Encoding
}

func (h *Handle) OutputFormat() OutputFormat {
	var rate C.long
	var channels C.int
	var encoding C.int
	C.mpg123_getformat(h.h, &rate, &channels, &encoding)

	return OutputFormat{
		Rate:     int(rate),
		Channels: int(channels),
		Encoding: Encoding(encoding),
	}
}

func (h *Handle) SetOutputFormat(f OutputFormat) error {
	C.mpg123_format_none(h.h)
	return mpgError(C.mpg123_format(h.h,
		C.long(f.Rate), C.int(f.Channels), C.int(f.Encoding)))

}

type MetaFlags int

const (
	MetaNewID3 MetaFlags = 0x1
	MetaNewICY           = 0x4
)

type ID3v2 struct {
	Version  uint8
	Title    string
	Artist   string
	Album    string
	Year     string
	Genre    string
	Comment  string
	Comments []ID3v2Text
	Text     []ID3v2Text
	Extra    []ID3v2Text
}

type ID3v2Text struct {
	Lang        string
	ID          string
	Description string
	Text        string
}

// Meta stores ID3v2 metadata, if any
type Meta struct {
	ID3v2 *ID3v2
}

func (h *Handle) MetaCheck() MetaFlags {
	return MetaFlags(C.mpg123_meta_check(h.h))
}

func convertID3v2String(s *C.mpg123_string) string {
	if s == nil {
		return ""
	}
	b := C.GoBytes(unsafe.Pointer(s.p), C.int(s.fill))
	return string(bytes.Join(bytes.Split(b, []byte{0x0}), []byte(" ")))

}

func convertID3v2TextArray(p *C.mpg123_text, size C.size_t) []ID3v2Text {
	if p == nil {
		return make([]ID3v2Text, 0)
	}
	res := make([]ID3v2Text, int(size))
	offset := unsafe.Sizeof(*p)
	for i := 0; i < len(res); i++ {
		t := (*C.mpg123_text)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + offset*uintptr(i)))
		if t == nil {
			continue
		}

		res[i] = ID3v2Text{
			Lang:        C.GoStringN(&t.lang[0], 3),
			ID:          C.GoStringN(&t.id[0], 4),
			Description: C.GoString(t.description.p),
			Text:        C.GoString(t.text.p),
		}

	}
	return res
}

func (h *Handle) MetaID3() (*ID3v2, error) {
	var v1 *C.mpg123_id3v1
	var v2 *C.mpg123_id3v2
	if err := C.mpg123_id3(h.h, &v1, &v2); err != C.MPG123_OK {
		return nil, mpgError(err)
	}
	if v2 != nil {
		return &ID3v2{
			Version:  uint8(v2.version),
			Title:    convertID3v2String(v2.title),
			Artist:   convertID3v2String(v2.artist),
			Album:    convertID3v2String(v2.album),
			Year:     convertID3v2String(v2.year),
			Genre:    convertID3v2String(v2.genre),
			Comment:  convertID3v2String(v2.comment),
			Comments: convertID3v2TextArray(v2.comment_list, v2.comments),
			Text:     convertID3v2TextArray(v2.text, v2.texts),
			Extra:    convertID3v2TextArray(v2.extra, v2.extras),
		}, nil
	}

	return nil, nil

}

func (h *Handle) SetFileSize(size int) {
	C.mpg123_set_filesize(h.h, C.off_t(size))
}

// Current stream offset
type Offset struct {
	Sample int
	Frame  int
	Byte   int
}

func (h *Handle) Offset() Offset {
	return Offset{
		Sample: int(C.mpg123_tell(h.h)),
		Frame:  int(C.mpg123_tellframe(h.h)),
		Byte:   int(C.mpg123_tell_stream(h.h)),
	}
}

//go:generate stringer -type=FrameMode

type FrameMode int

const (
	ModeStereo FrameMode = iota
	ModeJoint
	ModeDual
	ModeMono
)

//go:generate stringer -type=MpegVersion

type MpegVersion int

const (
	Mpeg10 MpegVersion = iota
	Mpeg20
	Mpeg25
)

type FrameFlags int

const (
	FrameCrc FrameFlags = 1 << iota
	FrameCopyright
	FramePrivate
	FrameOriginal
)

//go:generate stringer -type=FrameVbrMode

type FrameVbrMode int

const (
	Cbr FrameVbrMode = iota
	Vbr
	Abr
)

type FrameInfo struct {
	Version   MpegVersion
	Layer     int
	Rate      int
	Mode      FrameMode
	ModeExt   int
	Framesize int
	Flags     FrameFlags
	Emphasis  int
	Bitrate   int
	AbrRate   int
	Vbr       FrameVbrMode
}

func (h *Handle) FrameInfo() FrameInfo {
	var info C.struct_mpg123_frameinfo
	if C.mpg123_info(h.h, &info) == C.MPG123_OK {
		return FrameInfo{
			Version:   MpegVersion(info.version),
			Layer:     int(info.layer),
			Rate:      int(info.rate),
			Mode:      FrameMode(info.mode),
			ModeExt:   int(info.mode_ext),
			Framesize: int(info.framesize),
			Flags:     FrameFlags(info.flags),
			Emphasis:  int(info.emphasis),
			Bitrate:   int(info.bitrate),
			AbrRate:   int(info.abr_rate),
			Vbr:       FrameVbrMode(info.vbr),
		}
	} else {
		return FrameInfo{}
	}
}
