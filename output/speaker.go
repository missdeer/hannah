package output

import (
	"time"
)

type SpeakerStatus struct {
	Position time.Duration
	Length   time.Duration
	Volume   float64
	Speed    float64
}

type ISpeaker interface {
	IsPaused() bool
	Play()
	Init(int, int)
	Shutdown()
	PauseResume()
	Backward()
	Forward()
	IncreaseVolume()
	DecreaseVolume()
	Slowdown()
	Speedup()
	Status() *SpeakerStatus
	IsNil() bool
}
