package bass

import "C"
import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	pluginHandles []int
)

type Speaker struct {
	volumeRate float64
	volumeBase float64
	speedRate  float64
	freqBase   float64
	handle     uint
	paused     bool
	done       chan struct{}
}

func NewSpeaker() *Speaker {
	return &Speaker{volumeRate: 100.0, speedRate: 100.0}
}

func (s *Speaker) Initialize() {
	Init()

	dirs, reg := pluginsPattern()
	for _, dir := range dirs {
		if stat, err := os.Stat(dir); os.IsNotExist(err) || !stat.IsDir() {
			continue
		}
		fi, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, f := range fi {
			if f.IsDir() {
				continue
			}
			if !reg.MatchString(f.Name()) {
				continue
			}
			fn := filepath.Join(dir, f.Name())
			h := PluginLoad(fn)
			pluginHandles = append(pluginHandles, h)
		}
	}
}

func (s *Speaker) Finalize() {
	ChannelStop(s.handle)
	Free()
	for _, h := range pluginHandles {
		PluginFree(h)
	}
	Free()
}

func (s *Speaker) UpdateURI(uri string, done chan struct{}) {
	s.done = done
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		// http/https
		s.handle = StreamCreateURL(uri, 0, nil, nil)
	} else {
		// local file system
		s.handle = StreamCreateFile(0, uri, 0, 0)
	}
	s.freqBase = float64(GetChanAttr(s.handle, BASS_ATTRIB_FREQ))
	s.volumeBase = float64(GetChanAttr(s.handle, BASS_ATTRIB_VOL))
}

func (s *Speaker) UpdateStream(sampleRate int, streamer interface{}, done chan struct{}) {
}

func (s *Speaker) IsPaused() bool {
	return s.paused
}

func (s *Speaker) Play() {
	ChannelPlay(s.handle, 0)
}

func (s *Speaker) PrePlay(sampleRete int, bufferSize int) {
}

func (s *Speaker) Shutdown() {
	ChannelStop(s.handle)
	StreamFree(s.handle)
}

func (s *Speaker) PauseResume() {
	s.paused = !s.paused
	if s.paused {
		ChannelPause(s.handle)
	} else {
		ChannelPlay(s.handle, 1)
	}
}

func (s *Speaker) Backward() {
	pos := s.getCurrentPosition()
	pos -= 5 * time.Second
	if pos < 0 {
		pos = 0
	}
	ChannelSetPosition(s.handle, BASS_POS_BYTE, ChannelSeconds2Bytes(s.handle, int(pos/time.Millisecond)))
}

func (s *Speaker) Forward() {
	pos := s.getCurrentPosition()
	pos += 5 * time.Second
	length := s.getSongLength()
	if pos > length {
		pos = length
	}
	ChannelSetPosition(s.handle, BASS_POS_BYTE, ChannelSeconds2Bytes(s.handle, int(pos/time.Millisecond)))
}

func (s *Speaker) IncreaseVolume() {
	s.volumeRate = s.volumeRate * 1.1
	if s.volumeRate > 400 {
		s.volumeRate = 400
	}
	ChannelSetAttribute(s.handle, BASS_ATTRIB_VOL, C.float(s.volumeRate*s.volumeBase/100))
}

func (s *Speaker) DecreaseVolume() {
	s.volumeRate = s.volumeRate * 0.9
	if s.volumeRate < 10 {
		s.volumeRate = 10
	}
	ChannelSetAttribute(s.handle, BASS_ATTRIB_VOL, C.float(s.volumeRate*s.volumeBase/100))
}

func (s *Speaker) Slowdown() {
	s.speedRate *= 1.0 - 0.0594631
	if s.speedRate < 10 {
		s.speedRate = 10
	}
	ChannelSetAttribute(s.handle, BASS_ATTRIB_FREQ, C.float(s.speedRate*s.freqBase/100))
}

func (s *Speaker) Speedup() {
	s.speedRate *= 1.0594631 // 加速一次频率变为原来的(2的1/12次方=1.0594631)倍，即使单调提高一个半音，减速时同理
	if s.speedRate > 400 {
		s.speedRate = 400
	}
	ChannelSetAttribute(s.handle, BASS_ATTRIB_FREQ, C.float(s.speedRate*s.freqBase/100))
}

func (s *Speaker) getCurrentPosition() time.Duration {
	posInBytes := ChannelGetPosition(s.handle, BASS_POS_BYTE)
	posInSeconds := ChannelBytes2Seconds(s.handle, posInBytes)
	currentPosition := int(posInSeconds * 1000)
	if currentPosition == -1000 {
		currentPosition = 0
	}
	return time.Duration(currentPosition) * time.Microsecond
}

func (s *Speaker) getSongLength() time.Duration {
	lengthBytes := ChannelGetLength(s.handle, BASS_POS_BYTE)
	lengthSeconds := ChannelBytes2Seconds(s.handle, lengthBytes)
	songLength := int(lengthSeconds * 1000)
	if songLength == -1000 {
		songLength = 0
	}
	return time.Duration(songLength) * time.Microsecond
}

func (s *Speaker) Status() (time.Duration, time.Duration, float64, float64) {
	return s.getCurrentPosition(),
		s.getSongLength(),
		s.volumeRate / 100,
		s.speedRate / 100
}

func (s *Speaker) IsNil() bool {
	return false
}
