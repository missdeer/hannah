package decode

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/faiface/beep"
	"github.com/pkg/errors"

	"github.com/missdeer/hannah/media/decode/faad"
)

const (
	faadNumChannels   = 2
	faadPrecision     = 2
	faadBytesPerFrame = faadNumChannels * faadPrecision
)

func FAADDecode(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "m4a")
		}
	}()

	const faadSampeRate = 44100
	format = beep.Format{
		SampleRate:  beep.SampleRate(faadSampeRate),
		NumChannels: faadNumChannels,
		Precision:   faadPrecision,
	}

	context := faad.FaadDecoderCreate(faadSampeRate, faadNumChannels, faadPrecision)
	if context == nil {
		return nil, beep.Format{}, errors.New("creating FAAD context failed")
	}

	return &faaddecoder{
		inData:  rc,
		format:  format,
		context: context,
	}, format, nil
}

type faaddecoder struct {
	inData  io.ReadCloser
	format  beep.Format
	context unsafe.Pointer
	pos     int
	length  int
	err     error
}

func (d *faaddecoder) Stream(samples [][2]float64) (n int, ok bool) {
	if d.err != nil {
		return 0, false
	}

	for i := range samples {
		var inData [2048]byte
		dn, err := io.ReadFull(d.inData, inData[:])
		if err != nil {
			d.err = errors.Wrap(err, "m4a")
			break
		}
		var frame [2048]byte
		var frameLen int
		if faad.FaadGetOneADTSFrame(inData[:], dn, frame[:], &frameLen) != 0 {
			break
		}
		var outData [4096]byte
		var outLen int
		faad.FaadDecodeFrame(d.context, frame[:], frameLen, outData[:], &outLen)
		if outLen == len(outData) {
			samples[i], _ = d.format.DecodeSigned(outData[:])
			d.pos += outLen
			n++
			ok = true
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			d.err = errors.Wrap(err, "m4a")
			break
		}
	}
	return n, ok
}

func (d *faaddecoder) Err() error {
	return d.err
}

func (d *faaddecoder) Len() int {
	return int(d.length / faadBytesPerFrame)
}

func (d *faaddecoder) Position() int {
	return d.pos / faadBytesPerFrame
}

func (d *faaddecoder) Seek(p int) error {
	if p < 0 || d.Len() < p {
		return fmt.Errorf("mp3: seek position %v out of range [%v, %v]", p, 0, d.Len())
	}

	return nil
}

func (d *faaddecoder) Close() error {
	err := d.inData.Close()
	if err != nil {
		return errors.Wrap(err, "m4a")
	}
	faad.FaadDecodeClose(d.context)
	return nil
}
