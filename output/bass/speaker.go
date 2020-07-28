package bass

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
	handle uint
	paused bool
	done   chan struct{}
}

func NewSpeaker() *Speaker {
	return &Speaker{}
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
	pos := ChannelGetPosition(s.handle, BASS_POS_BYTE)
	pos -= 10
	if pos < 0 {
		pos = 0
	}
	ChannelSetPosition(s.handle, BASS_POS_BYTE, pos)
}

func (s *Speaker) Forward() {
	pos := ChannelGetPosition(s.handle, BASS_POS_BYTE)
	pos += 10
	length := ChannelGetLength(s.handle, BASS_POS_BYTE)
	if pos > length {
		pos = length
	}
	ChannelSetPosition(s.handle, BASS_POS_BYTE, pos)
}

func (s *Speaker) IncreaseVolume() {
	vol := GetChanVol(s.handle)
	SetChanVol(s.handle, uint(float64(vol)*1.1))
}

func (s *Speaker) DecreaseVolume() {
	vol := GetChanVol(s.handle)
	SetChanVol(s.handle, uint(float64(vol)*0.9))
}

func (s *Speaker) Slowdown() {
}

func (s *Speaker) Speedup() {
}

func (s *Speaker) Status() (time.Duration, time.Duration, float64, float64) {
	return time.Duration(ChannelGetPosition(s.handle, BASS_POS_BYTE)),
		time.Duration(ChannelGetLength(s.handle, BASS_POS_BYTE)),
		float64(GetVolume()),
		0
}

func (s *Speaker) IsNil() bool {
	return false
}
