package output

import (
	"strings"
	"time"

	"github.com/missdeer/hannah/output/bass"
)

type ISpeaker interface {
	Initialize()
	Finalize()
	IsPaused() bool
	Play()
	PrePlay(int, int)
	Shutdown()
	PauseResume()
	Backward()
	Forward()
	IncreaseVolume()
	DecreaseVolume()
	Slowdown()
	Speedup()
	Status() (time.Duration, time.Duration, float64, float64)
	IsNil() bool
	UpdateURI(string, chan struct{})
	UpdateStream(int, interface{}, chan struct{})
}

func NewSpeaker(engine string) ISpeaker {
	switch strings.ToLower(engine) {
	case "bass":
		return bass.NewSpeaker()
	}
	return nil
}
