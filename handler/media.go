package handler

import (
	"errors"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"

	"github.com/missdeer/hannah/input"
	"github.com/missdeer/hannah/output"
)

var (
	ShouldQuit   = errors.New("should quit application now")
	PreviousSong = errors.New("Play previous song")
	NextSong     = errors.New("Play next song")
	ap           *output.AudioPanel
	screen       tcell.Screen
	tcellEvents  = make(chan tcell.Event)
)

func PlayMedia(uri string, index int, total int, artist string, title string) error {
	r, err := input.OpenSource(uri)
	if err != nil {
		return err
	}
	defer r.Close()

	decoder := getDecoder(uri)
	streamer, format, err := decoder(r)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	defer func() {
		speaker.Clear()
		speaker.Close()
	}()

	done := make(chan struct{})
	if ap == nil {
		ap = output.NewAudioPanel(format.SampleRate, streamer, uri, index, total, artist, title, done)

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
		ap.Update(format.SampleRate, streamer, uri, index, total, artist, title, done)
	}

	screen.Clear()
	ap.Draw(screen)
	screen.Show()

	ap.Play()

	seconds := time.Tick(time.Second)
	for {
		select {
		case event := <-tcellEvents:
			changed, action := ap.Handle(event)
			switch action {
			case output.HandleActionQUIT:
				screen.Fini()
				return ShouldQuit
			case output.HandleActionNEXT:
				return NextSong
			case output.HandleActionPREVIOUS:
				return PreviousSong
			default:
				if changed {
					screen.Clear()
					ap.Draw(screen)
					screen.Show()
				}
			}
		case <-seconds:
			if !ap.IsPaused() {
				screen.Clear()
				ap.Draw(screen)
				screen.Show()
			}
		case <-done:
			return NextSong
		}
	}
	return nil
}
