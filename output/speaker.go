package output

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

type SpeakerStatus struct {
	Position time.Duration
	Length   time.Duration
	Volume   float64
	Speed    float64
}

func NewSpeaker(sampleRate beep.SampleRate, streamer beep.StreamSeeker, done chan struct{}) *Speaker {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &Speaker{
		sampleRate: sampleRate,
		streamer:   streamer,
		ctrl:       ctrl,
		resampler:  resampler,
		volume:     volume,
		done:       done,
	}
}

func (s *Speaker) Update(sampleRate beep.SampleRate, streamer beep.StreamSeeker, done chan struct{}) {
	s.sampleRate = sampleRate
	s.streamer = streamer
	s.ctrl = &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
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

func (s *Speaker) Initialize(sampleRate beep.SampleRate, bufferSize int) {
	speaker.Init(sampleRate, bufferSize)
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

func (s *Speaker) Status() *SpeakerStatus {
	speaker.Lock()
	defer speaker.Unlock()
	return &SpeakerStatus{
		Position: s.sampleRate.D(s.streamer.Position()),
		Length:   s.sampleRate.D(s.streamer.Len()),
		Volume:   s.volume.Volume,
		Speed:    s.resampler.Ratio(),
	}
}
