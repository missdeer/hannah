package faad

// #cgo pkg-config: faad2
// #cgo windows LDFLAGS: -lws2_32
// #include "neaacdec.h"
// #include "faacdecoder.h"
import "C"

import (
	"unsafe"
)

// FaadGetOneADTSFrame get one ADTS frame
func FaadGetOneADTSFrame(inData []byte, inLen int, outPCM []byte, outLen *int) int {
	return int(C.get_one_ADTS_frame((*C.uchar)(unsafe.Pointer(&inData[0])), C.size_t(inLen), (*C.uchar)(unsafe.Pointer(&outPCM[0])), (*C.size_t)(unsafe.Pointer(outLen))))
}

// FaadDecoderCreate create and return FAAD decoder context
func FaadDecoderCreate(sampleRate int, channels int, bitRate int) unsafe.Pointer {
	return C.faad_decoder_create(C.int(sampleRate), C.int(channels), C.int(bitRate))
}

// FaadDecodeFrame decode FAAC data and return PCM data, outPCM must large than 4096 bytes
func FaadDecodeFrame(param unsafe.Pointer, inData []byte, inLen int, outPCM []byte, outLen *int) {
	C.faad_decode_frame(param, (*C.uchar)(unsafe.Pointer(&inData[0])), C.int(inLen), (*C.uchar)(unsafe.Pointer(&outPCM[0])), (*C.uint)(unsafe.Pointer(outLen)))
}

// FaadDecodeClose close FAAD decoder context
func FaadDecodeClose(param unsafe.Pointer) {
	C.faad_decode_close(param)
}
