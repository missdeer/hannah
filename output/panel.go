package output

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
)

const (
	HandleActionQUIT = iota
	HandleActionPREVIOUS
	HandleActionNEXT
	HandleActionNOP
)

func drawTextLine(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for _, r := range s {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

type AudioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	mediaURI   string
	index      int
	total      int
	done       chan struct{}
}

func NewAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker, uri string, index int, total int, done chan struct{}) *AudioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &AudioPanel{
		sampleRate,
		streamer,
		ctrl,
		resampler,
		volume,
		uri,
		index,
		total,
		done,
	}
}

func (ap *AudioPanel) Play() {
	speaker.Play(beep.Seq(ap.volume, beep.Callback(func() {
		ap.done <- struct{}{}
	})))
}

func (ap *AudioPanel) Draw(screen tcell.Screen) {
	mainStyle := tcell.StyleDefault.
		Background(tcell.NewHexColor(0x473437)).
		Foreground(tcell.NewHexColor(0xD7D8A2))
	statusStyle := mainStyle.
		Foreground(tcell.NewHexColor(0xDDC074)).
		Bold(true)

	screen.Fill(' ', mainStyle)

	drawTextLine(screen, 0, 0, "Welcome to Hannah, the simplest music player!", mainStyle)
	drawTextLine(screen, 0, 1, "Press [ESC] to quit.", mainStyle)
	drawTextLine(screen, 0, 2, "Press [SPACE] to pause/resume.", mainStyle)
	drawTextLine(screen, 0, 3, "Use keys in (?/?) to turn the buttons.", mainStyle)

	speaker.Lock()
	position := ap.sampleRate.D(ap.streamer.Position())
	length := ap.sampleRate.D(ap.streamer.Len())
	volume := ap.volume.Volume
	speed := ap.resampler.Ratio()
	speaker.Unlock()

	positionStatus := fmt.Sprintf("%v / %v", position.Round(time.Second), length.Round(time.Second))
	volumeStatus := fmt.Sprintf("%.1f", volume)
	speedStatus := fmt.Sprintf("%.3fx", speed)

	s := fmt.Sprintf("Media [%d/%d] (P/N):", ap.index, ap.total)
	drawTextLine(screen, 0, 5, s, mainStyle)
	drawTextLine(screen, len(s), 5, ap.mediaURI, statusStyle)

	drawTextLine(screen, 0, 6, "Position"+strings.Repeat(" ", len(s)-len(`Position`)-len(`(Q/W):`))+"(Q/W):", mainStyle)
	drawTextLine(screen, len(s), 6, positionStatus, statusStyle)

	drawTextLine(screen, 0, 7, "Volume"+strings.Repeat(" ", len(s)-len(`Volume`)-len(`(A/S):`))+"(A/S):", mainStyle)
	drawTextLine(screen, len(s), 7, volumeStatus, statusStyle)

	drawTextLine(screen, 0, 8, "Speed"+strings.Repeat(" ", len(s)-len(`Speed`)-len(`(Z/X):`))+"(Z/X):", mainStyle)
	drawTextLine(screen, len(s), 8, speedStatus, statusStyle)
}

func (ap *AudioPanel) Handle(event tcell.Event) (changed bool, action int) {
	switch event := event.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyESC {
			return false, HandleActionQUIT
		}

		if event.Key() != tcell.KeyRune {
			return false, HandleActionNOP
		}

		switch unicode.ToLower(event.Rune()) {
		case ' ':
			speaker.Lock()
			ap.ctrl.Paused = !ap.ctrl.Paused
			speaker.Unlock()
			return false, HandleActionNOP

		case 'q', 'w':
			speaker.Lock()
			newPos := ap.streamer.Position()
			if event.Rune() == 'q' {
				newPos -= ap.sampleRate.N(time.Second)
			}
			if event.Rune() == 'w' {
				newPos += ap.sampleRate.N(time.Second)
			}
			if newPos < 0 {
				newPos = 0
			}
			if newPos >= ap.streamer.Len() {
				newPos = ap.streamer.Len() - 1
			}
			if err := ap.streamer.Seek(newPos); err != nil {
				log.Fatal(err)
			}
			speaker.Unlock()
			return true, HandleActionNOP

		case 'a':
			speaker.Lock()
			ap.volume.Volume -= 0.1
			speaker.Unlock()
			return true, HandleActionNOP

		case 's':
			speaker.Lock()
			ap.volume.Volume += 0.1
			speaker.Unlock()
			return true, HandleActionNOP

		case 'z':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 15 / 16)
			speaker.Unlock()
			return true, HandleActionNOP

		case 'x':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 16 / 15)
			speaker.Unlock()
			return true, HandleActionNOP

		case 'p':
			speaker.Clear()
			speaker.Close()
			return true, HandleActionPREVIOUS

		case 'n':
			speaker.Clear()
			speaker.Close()
			return true, HandleActionNEXT
		}
	}
	return false, HandleActionNOP
}
