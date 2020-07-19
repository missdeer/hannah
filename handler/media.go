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
)

func PlayMedia(uri string, index int, total int) error {
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

	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	err = screen.Init()
	if err != nil {
		return err
	}
	defer screen.Fini()

	done := make(chan struct{})
	ap := output.NewAudioPanel(format.SampleRate, streamer, uri, index, total, done)

	screen.Clear()
	ap.Draw(screen)
	screen.Show()

	ap.Play()

	seconds := time.Tick(time.Second)
	events := make(chan tcell.Event)
	go func() {
		for {
			events <- screen.PollEvent()
		}
	}()

	for {
		select {
		case event := <-events:
			changed, action := ap.Handle(event)
			switch action {
			case output.HandleActionQUIT:
				return ShouldQuit
			case output.HandleActionNEXT:
				return NextSong
			case output.HandleActionPREVIOUS:
				return PreviousSong
			}
			if changed {
				screen.Clear()
				ap.Draw(screen)
				screen.Show()
			}
		case <-seconds:
			screen.Clear()
			ap.Draw(screen)
			screen.Show()
		case <-done:
			return NextSong
		}
	}
	return nil
}
