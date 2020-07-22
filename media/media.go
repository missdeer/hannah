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
	screenPanel          *output.ScreenPanel
	audioSpeaker         *output.Speaker
	screen               tcell.Screen
	tcellEvents          = make(chan tcell.Event)
)

func PlayMedia(uri string, index int, total int, artist string, title string) error {
	if screenPanel != nil {
		screenPanel.SetMessage(fmt.Sprintf("Loading %s ...", uri))
		status := audioSpeaker.Status()
		screenPanel.Draw(screen, status.Position, status.Length, status.Volume, status.Speed)
	}
	r, err := input.OpenSource(uri)
	if err != nil {
		return err
	}
	defer r.Close()

	if screenPanel != nil {
		screenPanel.SetMessage(fmt.Sprintf("Decoding %s ...", uri))
		status := audioSpeaker.Status()
		screenPanel.Draw(screen, status.Position, status.Length, status.Volume, status.Speed)
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

	if screenPanel != nil {
		screenPanel.SetMessage("Initializing speaker...")
		status := audioSpeaker.Status()
		screenPanel.Draw(screen, status.Position, status.Length, status.Volume, status.Speed)
	}

	done := make(chan struct{})

	if audioSpeaker == nil {
		audioSpeaker = output.NewSpeaker(format.SampleRate, streamer, done)
		audioSpeaker.Initialize(format.SampleRate, format.SampleRate.N(time.Second/10))
	} else {
		audioSpeaker.Update(format.SampleRate, streamer, done)
		audioSpeaker.Initialize(format.SampleRate, format.SampleRate.N(time.Second/10))
	}
	defer audioSpeaker.Shutdown()

	if screenPanel == nil {
		screenPanel = output.NewScreenPanel(uri, index, total, artist, title)

		screen, err = tcell.NewScreen()
		if err != nil {
			return err
		}
		err = screen.Init()
		if err != nil {
			return err
		}

		go func() {
			for {
				tcellEvents <- screen.PollEvent()
			}
		}()
	} else {
		screenPanel.Update(uri, index, total, artist, title)
	}

	screenPanel.SetMessage("")
	status := audioSpeaker.Status()
	screenPanel.Draw(screen, status.Position, status.Length, status.Volume, status.Speed)

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
				screenPanel.Draw(screen, status.Position, status.Length, status.Volume, status.Speed)
			}
		case <-seconds:
			if !audioSpeaker.IsPaused() {
				status := audioSpeaker.Status()
				screenPanel.Draw(screen, status.Position, status.Length, status.Volume, status.Speed)
			}
		case <-done:
			return NextSong
		}
	}
	return nil
}

func Finalize() {
	if screen != nil {
		screen.Fini()
	}
}
