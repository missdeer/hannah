package media

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"

	"github.com/missdeer/hannah/bass"
	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/input"
	"github.com/missdeer/hannah/media/decode"
	"github.com/missdeer/hannah/output"
)

type play func(string, int, int, string, string) error

var (
	ShouldQuit           = errors.New("should quit application now")
	PreviousSong         = errors.New("play previous song")
	NextSong             = errors.New("play next song")
	UnsupportedMediaType = errors.New("unsupported media type")
	screenPanel          *output.ScreenPanel
	audioSpeaker         output.ISpeaker
	tcellEvents          chan tcell.Event
	PlayMedia            play = builtinPlayMedia
)

func Initialize() error {
	switch strings.ToLower(config.Engine) {
	case "builtin":
		PlayMedia = builtinPlayMedia
	case "bass":
		bass.Init()
		PlayMedia = bassPlayMedia
	}
	audioSpeaker = output.NewSpeaker(config.Engine)

	screenPanel = output.NewScreenPanel()
	if err := screenPanel.Initialize(); err != nil {
		return err
	}

	go func() {
		tcellEvents = make(chan tcell.Event)
		defer func() {
			close(tcellEvents)
		}()
		for ; !screenPanel.Quit(); {
			tcellEvents <- screenPanel.PollScreenEvent()
		}
	}()
	return nil
}

func bassPlayMedia(uri string, index int, total int, artist string, title string) error {
	return nil
}

func builtinPlayMedia(uri string, index int, total int, artist string, title string) error {
	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Loading %s ...", uri))
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
	}

	r, err := input.OpenSource(uri)
	if err != nil {
		return err
	}
	defer r.Close()

	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Decoding %s ...", uri))
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
	}

	decoder := decode.GetBuiltinDecoder(uri)
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
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
	}

	audioSpeaker.Init(int(format.SampleRate), format.SampleRate.N(time.Second/10))
	defer audioSpeaker.Shutdown()

	done := make(chan struct{})
	audioSpeaker.UpdateStream(int(format.SampleRate), streamer, done)

	screenPanel.Update(uri, index, total, artist, title)

	screenPanel.SetMessage("")
	pos, length, vol, speed := audioSpeaker.Status()
	screenPanel.Draw(pos, length, vol, speed)

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
				pos, length, vol, speed := audioSpeaker.Status()
				screenPanel.Draw(pos, length, vol, speed)
			}
		case <-seconds:
			if !audioSpeaker.IsPaused() {
				pos, length, vol, speed := audioSpeaker.Status()
				screenPanel.Draw(pos, length, vol, speed)
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
