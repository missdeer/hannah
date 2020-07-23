package mpg123

import (
	"errors"
	"io"
	"io/ioutil"
)

// Reader implements on-the-fly mp3 decoding
// for an io.Reader
type Reader struct {
	input        io.Reader
	h            *Handle
	totalBytes   int64
	feedBuf      []byte
	nextOffset   int64
	totalRead    int64
	bytesSinceOk int
	maxBadBytes  int
	needMore     bool
	offset       Offset
	frameInfo    FrameInfo
	outputFormat OutputFormat
	meta         Meta
}

type ReaderConfig struct {
	// Output sample format of the decoded audio
	OutputFormat *OutputFormat
	// Internal buffer size
	BufferSize int
	// Total length in byte
	TotalBytes int64
}

var DefaultConfig = ReaderConfig{
	OutputFormat: nil,
	BufferSize:   16 * 1024,
}

// NewReaderConfig returns a new Reader
// configured with supplied config
func NewReaderConfig(r io.Reader, config ReaderConfig) *Reader {
	h, _ := NewDefaultHandle()

	if config.OutputFormat != nil {
		h.SetOutputFormat(*config.OutputFormat)
	}

	if config.BufferSize == 0 {
		config.BufferSize = DefaultConfig.BufferSize
	}

	h.OpenFeed()
	return &Reader{
		input:       r,
		totalBytes:  config.TotalBytes,
		h:           h,
		maxBadBytes: 4096,
		feedBuf:     make([]byte, config.BufferSize),
	}
}

// NewReader returns a new Reader
// configured with the supplied config
func NewReader(r io.Reader) *Reader {
	return NewReaderConfig(r, ReaderConfig{})
}

// TotalBytes returns bytes count had been read
func (r *Reader) TotalBytes() int64 {
	return r.totalBytes
}

// Offset returns current stream offset
func (r *Reader) Offset() Offset {
	return r.offset
}

// FrameInfo returns the most recent frame info
func (r *Reader) FrameInfo() FrameInfo {
	return r.frameInfo
}

// OutputFormat returns the decoder output format
func (r *Reader) OutputFormat() OutputFormat {
	return r.outputFormat
}

// Meta returns stream metadata
func (r *Reader) Meta() Meta {
	return r.meta
}

// Read implements io.Reader, reading decoded
// samples from its underlying Reader.
// The sample format depends on the configuration.
// Read(nil) allows to find the beginning the stream
// without  consuming any  samples
func (r *Reader) Read(buf []byte) (int, error) {
	defer func() {
		r.offset = r.h.Offset()
		r.frameInfo = r.h.FrameInfo()

		f := r.h.MetaCheck()
		switch {
		case f&MetaNewID3 != 0:
			id3v2, err := r.h.MetaID3()
			if id3v2 != nil && err == nil {
				r.meta.ID3v2 = id3v2
			}
		}

	}()
	if r.nextOffset > r.totalRead {
		n, err := io.CopyN(ioutil.Discard, r.input, r.nextOffset-r.totalRead)
		r.totalRead += n
		if err != nil {
			return 0, err
		}
	}
	for r.bytesSinceOk < r.maxBadBytes {
		var feed []byte
		if r.needMore {
			r.needMore = false
			feedLen, err := r.input.Read(r.feedBuf)
			r.totalRead += int64(feedLen)
			r.nextOffset = r.totalRead
			if feedLen == 0 && err != nil {
				return 0, err
			}
			feed = r.feedBuf[:feedLen]
			r.bytesSinceOk += feedLen
		}

		switch n, err := r.h.Decode(feed, buf); err {
		case ErrNewFormat:
			r.outputFormat = r.h.OutputFormat()
			r.bytesSinceOk = 0
			if len(buf) == 0 {
				return n, nil
			}
		case ErrNeedMore:
			r.needMore = true
			if n > 0 {
				r.bytesSinceOk = 0
				return n, nil
			}
		case ErrDone:
			return n, io.EOF
		default:
			r.bytesSinceOk = 0
			return n, nil

		}

	}
	r.bytesSinceOk = 0
	return 0, errors.New("No valid data found")
}

// Seek sets the offset in samples of the next Read.
// If underlying reader is not io.Seeker only io.SeekStart and io.SeekCurrent
// with positive offsets are supported
func (r *Reader) Seek(sampleOffset int64, whence int) (int64, error) {
	newOffset, inputBytesOffset, err := r.h.FeedSeek(sampleOffset, whence)
	if err != nil {
		return 0, err
	}
	switch ir := r.input.(type) {
	case io.Seeker:
		if _, err := ir.Seek(int64(inputBytesOffset), io.SeekStart); err != nil {
			return 0, err
		}
	default:
		if inputBytesOffset < r.totalRead {
			return 0, errors.New("Seeking back is not supported")
		}
		r.nextOffset = inputBytesOffset
	}
	return newOffset, nil
}
