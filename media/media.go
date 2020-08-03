package media

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/input"
	"github.com/missdeer/hannah/media/decode"
	"github.com/missdeer/hannah/output"
	"github.com/missdeer/hannah/output/bass"
	"github.com/missdeer/hannah/output/beep"
)

type play func(string, int, int, string, string) error
type supportedFileType func(string) bool

var (
	ShouldQuit           = errors.New("should quit application now")
	PreviousSong         = errors.New("play previous song")
	NextSong             = errors.New("play next song")
	UnsupportedMediaType = errors.New("unsupported media type")
	screenPanel          *output.ScreenPanel
	audioSpeaker         output.ISpeaker
	tcellEvents          chan tcell.Event
	PlayMedia            play              = builtinPlayMedia
	IsSupportedFileType  supportedFileType = beep.SupportedFileType
)

func Initialize(screenPanelEnabled bool) error {
	audioSpeaker = output.NewSpeaker(config.Engine)
	audioSpeaker.Initialize()

	switch strings.ToLower(config.Engine) {
	case "builtin":
		PlayMedia = builtinPlayMedia
		IsSupportedFileType = beep.SupportedFileType
	case "bass":
		PlayMedia = bassPlayMedia
		IsSupportedFileType = bass.SupportedFileType
	}
	if screenPanelEnabled {
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
	}
	return nil
}

func playMedia(uri string, index int, total int, artist string, title string, done chan struct{}) error {
	screenPanel.Update(uri, index, total, artist, title)

	screenPanel.SetMessage("")
	pos, length, vol, speed := audioSpeaker.Status()
	screenPanel.Draw(pos, length, vol, speed)

	audioSpeaker.Play()

	seconds := time.Tick(time.Second)
	lastPos := pos
	tickCount := 0
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
				if lastPos != pos {
					tickCount = 0
					lastPos = pos
				}
			}
		case <-seconds:
			if !audioSpeaker.IsPaused() {
				pos, length, vol, speed := audioSpeaker.Status()
				screenPanel.Draw(pos, length, vol, speed)
				if lastPos != pos {
					tickCount = 0
					lastPos = pos
				} else {
					tickCount++
					if tickCount > 10 { // doesn't play in 10 seconds, switch to next song
						return NextSong
					}
				}
			}
		case <-done:
			return NextSong
		}
	}
	return nil
}

func bassPlayMedia(uri string, index int, total int, artist string, title string) error {
	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Loading %s ...", uri))
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
	}

	defer audioSpeaker.Shutdown()

	done := make(chan struct{})
	audioSpeaker.UpdateURI(uri, done)

	return playMedia(uri, index, total, artist, title, done)
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

	audioSpeaker.PrePlay(int(format.SampleRate), format.SampleRate.N(time.Second/10))
	defer audioSpeaker.Shutdown()

	done := make(chan struct{})
	audioSpeaker.UpdateStream(int(format.SampleRate), streamer, done)

	return playMedia(uri, index, total, artist, title, done)
}

func Finalize(screenPanelEnabled bool) {
	if screenPanelEnabled {
		screenPanel.Finalize()
	}
	audioSpeaker.Finalize()
}
