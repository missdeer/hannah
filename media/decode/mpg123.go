package decode

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/faiface/beep"
	"github.com/pkg/errors"

	"github.com/missdeer/hannah/media/decode/mpg123"
)

const (
	gomp3NumChannels   = 2
	gomp3Precision     = 2
	gomp3BytesPerFrame = gomp3NumChannels * gomp3Precision
)

// Decode takes a ReadCloser containing audio data in MP3 format and returns a StreamSeekCloser,
// which streams that audio. The Seek method will panic if rc is not io.Seeker.
//
// Do not close the supplied ReadSeekCloser, instead, use the Close method of the returned
// StreamSeekCloser when you want to release the resources.
func Mpg123Decode(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "mp3")
		}
	}()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
		return nil, beep.Format{}, err
	}
	r := mpg123.NewReaderConfig(bytes.NewReader(b), mpg123.ReaderConfig{
		OutputFormat: &mpg123.OutputFormat{
			Channels: 2,
			Rate:     44100,
			Encoding: mpg123.EncodingInt16,
		},
	})

	r.Read(nil)
	format = beep.Format{
		SampleRate:  beep.SampleRate(r.OutputFormat().Rate),
		NumChannels: gomp3NumChannels,
		Precision:   gomp3Precision,
	}
	d, err := ioutil.ReadAll(r)
	length := 0

	if err == nil {
		length = len(d)
	}
	r = mpg123.NewReaderConfig(bytes.NewReader(b), mpg123.ReaderConfig{
		OutputFormat: &mpg123.OutputFormat{
			Channels: 2,
			Rate:     44100,
			Encoding: mpg123.EncodingInt16,
		},
	})
	return &decoder{rc, r, format, 0, length, nil}, format, nil
}

type decoder struct {
	closer io.Closer
	r      *mpg123.Reader
	f      beep.Format
	pos    int
	length int
	err    error
}

func (d *decoder) Stream(samples [][2]float64) (n int, ok bool) {
	if d.err != nil {
		return 0, false
	}
	var tmp [gomp3BytesPerFrame]byte
	for i := range samples {
		dn, err := d.r.Read(tmp[:])
		if dn == len(tmp) {
			samples[i], _ = d.f.DecodeSigned(tmp[:])
			d.pos += dn
			n++
			ok = true
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			d.err = errors.Wrap(err, "mp3")
			break
		}
	}
	return n, ok
}

func (d *decoder) Err() error {
	return d.err
}

func (d *decoder) Len() int {
	return int(d.length / gomp3BytesPerFrame)
}

func (d *decoder) Position() int {
	return d.pos / gomp3BytesPerFrame
}

func (d *decoder) Seek(p int) error {
	if p < 0 || d.Len() < p {
		return fmt.Errorf("mp3: seek position %v out of range [%v, %v]", p, 0, d.Len())
	}
	_, err := d.r.Seek(int64(p)*gomp3BytesPerFrame, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "mp3")
	}
	d.pos = p * gomp3BytesPerFrame
	return nil
}

func (d *decoder) Close() error {
	err := d.closer.Close()
	if err != nil {
		return errors.Wrap(err, "mp3")
	}
	return nil
}
