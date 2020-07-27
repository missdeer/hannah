package beep

import (
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

type Speaker struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	done       chan struct{}
}

func NewSpeaker() *Speaker {
	return &Speaker{}
}

func (s *Speaker) Initialize() {

}

func (s *Speaker) Finalize() {

}

func (s *Speaker) UpdateURI(uri string, done chan struct{}) {

}

func (s *Speaker) UpdateStream(sampleRate int, streamer interface{}, done chan struct{}) {
	s.sampleRate = beep.SampleRate(sampleRate)
	s.streamer = streamer.(beep.StreamSeeker)
	s.ctrl = &beep.Ctrl{Streamer: beep.Loop(1, s.streamer)}
	s.resampler = beep.ResampleRatio(4, 1, s.ctrl)
	s.volume = &effects.Volume{Streamer: s.resampler, Base: 2}
	s.done = done
}

func (s *Speaker) IsPaused() bool {
	speaker.Lock()
	defer speaker.Unlock()
	return s.ctrl.Paused
}

func (s *Speaker) Play() {
	speaker.Play(beep.Seq(s.volume, beep.Callback(func() {
		s.done <- struct{}{}
	})))
}

func (s *Speaker) PrePlay(sampleRete int, bufferSize int) {
	speaker.Init(beep.SampleRate(sampleRete), bufferSize)
}

func (s *Speaker) Shutdown() {
	speaker.Clear()
	speaker.Close()
}

func (s *Speaker) PauseResume() {
	speaker.Lock()
	s.ctrl.Paused = !s.ctrl.Paused
	speaker.Unlock()
}

func (s *Speaker) Backward() {
	speaker.Lock()
	if s.streamer.Len() > 0 {
		newPos := s.streamer.Position()
		newPos -= s.sampleRate.N(time.Second)
		if newPos < 0 {
			newPos = 0
		}
		if err := s.streamer.Seek(newPos); err != nil {
			log.Fatal(err)
		}
	}
	speaker.Unlock()
}

func (s *Speaker) Forward() {
	speaker.Lock()
	if s.streamer.Len() > 0 {
		newPos := s.streamer.Position()
		newPos += s.sampleRate.N(time.Second)
		if newPos >= s.streamer.Len() {
			newPos = s.streamer.Len() - 1
		}
		if err := s.streamer.Seek(newPos); err != nil {
			log.Fatal(err)
		}
	}
	speaker.Unlock()
}

func (s *Speaker) IncreaseVolume() {
	speaker.Lock()
	s.volume.Volume += 0.1
	speaker.Unlock()
}

func (s *Speaker) DecreaseVolume() {
	speaker.Lock()
	s.volume.Volume -= 0.1
	speaker.Unlock()
}

func (s *Speaker) Slowdown() {
	speaker.Lock()
	s.resampler.SetRatio(s.resampler.Ratio() * 15 / 16)
	speaker.Unlock()
}

func (s *Speaker) Speedup() {
	speaker.Lock()
	s.resampler.SetRatio(s.resampler.Ratio() * 16 / 15)
	speaker.Unlock()
}

func (s *Speaker) Status() (time.Duration, time.Duration, float64, float64) {
	speaker.Lock()
	defer speaker.Unlock()
	return s.sampleRate.D(s.streamer.Position()),
		s.sampleRate.D(s.streamer.Len()),
		s.volume.Volume,
		s.resampler.Ratio()
}
func (s *Speaker) IsNil() bool {
	return s.streamer == nil
}
