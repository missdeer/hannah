package output

import (
	"strings"
	"time"

	"github.com/missdeer/hannah/output/beep"
)

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
	Status() (time.Duration, time.Duration, float64, float64)
	IsNil() bool
	UpdateURI(int, string, chan struct{})
	UpdateStream(int, interface{}, chan struct{})
}

func NewSpeaker(engine string) ISpeaker {
	switch strings.ToLower(engine) {
	case "builtin":
		return beep.NewSpeaker()
	case "bass":
		return nil
	}
	return nil
}
