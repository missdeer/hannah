package bass

import (
	"strings"
	"time"

	"github.com/faiface/beep"
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
	ps := plugins()
	for _, plugin := range ps {
		h := PluginLoad(plugin)
		pluginHandles = append(pluginHandles, h)
	}
}

func (s *Speaker) Finalize() {
	for _, h := range pluginHandles {
		PluginFree(h)
	}
	Free()
}

func (s *Speaker) UpdateURI(sampleRate int, uri string, done chan struct{}) {
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

func (s *Speaker) InitializeSpeaker(sampleRate beep.SampleRate, bufferSize int) {
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
}

func (s *Speaker) Forward() {
}

func (s *Speaker) IncreaseVolume() {
}

func (s *Speaker) DecreaseVolume() {
}

func (s *Speaker) Slowdown() {
}

func (s *Speaker) Speedup() {
}

func (s *Speaker) Status() (time.Duration, time.Duration, float64, float64) {
	return 0, 0, 0, 0
}

func (s *Speaker) IsNil() bool {
	return false
}
