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
	"github.com/missdeer/hannah/provider"
)

var (
	ShouldQuit           = errors.New("should quit application now")
	PreviousSong         = errors.New("play previous song")
	NextSong             = errors.New("play next song")
	UnsupportedMediaType = errors.New("unsupported media type")
	screenPanel          *output.ScreenPanel
	audioSpeaker         output.ISpeaker
	tcellEvents          chan tcell.Event
	PlayMedia            = builtinPlayMedia
	IsSupportedFileType  = beep.SupportedFileType
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

func playMedia(song provider.Song, index int, total int, done chan struct{}) error {
	screenPanel.Update(song, index, total)

	screenPanel.SetMessage("")
	pos, length, vol, speed := audioSpeaker.Status()
	screenPanel.Draw(pos, length, vol, speed)

	audioSpeaker.Play()

	downloaded := make(chan struct{})
	saved := make(chan struct{})
	seconds := time.Tick(time.Second)
	lastPos := pos
	tickCount := 0

	updateStatus := func() bool {
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
		if lastPos != pos {
			tickCount = 0
			lastPos = pos
			return true
		}
		return false
	}
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
			case output.HandleActionM3U:
				go insertToM3U(song, saved)
			case output.HandleActionDownload:
				go downloadSong(song, downloaded)
			default:
			}
			if changed {
				updateStatus()
			}
		case <-seconds:
			if !audioSpeaker.IsPaused() {
				if !updateStatus() {
					tickCount++
					if tickCount > 10 { // doesn't play in 10 seconds, switch to next song
						return NextSong
					}
				}
			}
		case <-done:
			return NextSong
		case <-downloaded:
			screenPanel.SetMessage(fmt.Sprintf("%s-%s%s is save to %s", song.Title, song.Artist, decode.GetExtName(song.URL), config.DownloadDir))
			updateStatus()
		case <-saved:
			screenPanel.SetMessage(fmt.Sprintf("%s://%s is appended to %s", song.Provider, song.ID, config.M3UFileName))
			updateStatus()
		}
	}
	return nil
}

func bassPlayMedia(song provider.Song, index int, total int) error {
	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Loading %s ...", song.URL))
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
	}

	defer audioSpeaker.Shutdown()

	done := make(chan struct{})
	audioSpeaker.UpdateURI(song.URL, done)

	return playMedia(song, index, total, done)
}

func builtinPlayMedia(song provider.Song, index int, total int) error {
	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Loading %s ...", song.URL))
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
	}

	r, err := input.OpenSource(song.URL)
	if err != nil {
		return err
	}
	defer r.Close()

	if !audioSpeaker.IsNil() {
		screenPanel.SetMessage(fmt.Sprintf("Decoding %s ...", song.URL))
		pos, length, vol, speed := audioSpeaker.Status()
		screenPanel.Draw(pos, length, vol, speed)
	}

	decoder := decode.GetBuiltinDecoder(song.URL)
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

	return playMedia(song, index, total, done)
}

func Finalize(screenPanelEnabled bool) {
	if screenPanelEnabled {
		screenPanel.Finalize()
	}
	audioSpeaker.Finalize()
}
