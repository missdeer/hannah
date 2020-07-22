package media

import (
	"errors"
	"fmt"
	"time"

	"github.com/gdamore/tcell"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/input"
	"github.com/missdeer/hannah/output"
)

var (
	ShouldQuit           = errors.New("should quit application now")
	PreviousSong         = errors.New("play previous song")
	NextSong             = errors.New("play next song")
	UnsupportedMediaType = errors.New("unsupported media type")
	screenPanel          = output.NewScreenPanel()
	audioSpeaker         = output.NewSpeaker()
	tcellEvents          = make(chan tcell.Event)
)

func PlayMedia(uri string, index int, total int, artist string, title string) error {
	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Loading %s ...", uri))
		status := audioSpeaker.Status()
		screenPanel.Draw(status.Position, status.Length, status.Volume, status.Speed)
	}

	r, err := input.OpenSource(uri)
	if err != nil {
		return err
	}
	defer r.Close()

	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Decoding %s ...", uri))
		status := audioSpeaker.Status()
		screenPanel.Draw(status.Position, status.Length, status.Volume, status.Speed)
	}

	decoder := getDecoder(uri)
	if decoder == nil {
		return UnsupportedMediaType
	}
	streamer, format, err := decoder(r)
	if err != nil {
		return err
	}
	defer streamer.Close()

	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage("Initializing speaker...")
		status := audioSpeaker.Status()
		screenPanel.Draw(status.Position, status.Length, status.Volume, status.Speed)
	}

	audioSpeaker.InitializeSpeaker(format.SampleRate, format.SampleRate.N(time.Second/10))
	defer audioSpeaker.Shutdown()

	done := make(chan struct{})
	audioSpeaker.Update(format.SampleRate, streamer, done)

	go func() {
		for {
			tcellEvents <- screenPanel.PollScreenEvent()
		}
	}()

	screenPanel.Update(uri, index, total, artist, title)

	screenPanel.SetMessage("")
	status := audioSpeaker.Status()
	screenPanel.Draw(status.Position, status.Length, status.Volume, status.Speed)

	audioSpeaker.Play()

	seconds := time.Tick(time.Second)
	for {
		select {
		case event := <-tcellEvents:
			changed, action := screenPanel.Handle(event)
			switch action {
			case output.HandleActionQUIT:
				return ShouldQuit
			case output.HandleActionNEXT:
				return NextSong
			case output.HandleActionPREVIOUS:
				return PreviousSong
			case output.HandleActionRepeat:
				config.Repeat = !config.Repeat
			case output.HandleActionShuffle:
				config.Shuffle = !config.Shuffle
			case output.HandleActionPauseResume:
				audioSpeaker.PauseResume()
			case output.HandleActionBackward:
				audioSpeaker.Backward()
			case output.HandleActionForward:
				audioSpeaker.Forward()
			case output.HandleActionDecreaseVolume:
				audioSpeaker.DecreaseVolume()
			case output.HandleActionIncreaseVolume:
				audioSpeaker.IncreaseVolume()
			case output.HandleActionSlowdown:
				audioSpeaker.Slowdown()
			case output.HandleActionSpeedup:
				audioSpeaker.Speedup()
			default:
			}
			if changed {
				status := audioSpeaker.Status()
				screenPanel.Draw(status.Position, status.Length, status.Volume, status.Speed)
			}
		case <-seconds:
			if !audioSpeaker.IsPaused() {
				status := audioSpeaker.Status()
				screenPanel.Draw(status.Position, status.Length, status.Volume, status.Speed)
			}
		case <-done:
			return NextSong
		}
	}
	return nil
}

func Finalize() {
	screenPanel.Finalize()
}
