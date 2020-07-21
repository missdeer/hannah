package output

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/util"
)

const (
	HandleActionQUIT = iota
	HandleActionPREVIOUS
	HandleActionNEXT
	HandleActionNOP
)

func drawTextLine(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	text := []rune(s)
	for _, r := range text {
		screen.SetContent(x, y, r, nil, style)
		if utf8.RuneLen(r) > 1 {
			x += 2
		} else {
			x++
		}
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
	artist     string
	title      string
	message    string
	done       chan struct{}
}

func NewAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker, uri string, index int, total int, artist string, title string, done chan struct{}) *AudioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &AudioPanel{
		sampleRate: sampleRate,
		streamer:   streamer,
		ctrl:       ctrl,
		resampler:  resampler,
		volume:     volume,
		mediaURI:   uri,
		index:      index,
		total:      total,
		artist:     artist,
		title:      title,
		done:       done,
	}
}

func (ap *AudioPanel) Update(sampleRate beep.SampleRate, streamer beep.StreamSeeker, uri string, index int, total int, artist string, title string, done chan struct{}) {
	ap.sampleRate = sampleRate
	ap.streamer = streamer
	ap.ctrl = &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	ap.resampler = beep.ResampleRatio(4, 1, ap.ctrl)
	ap.volume = &effects.Volume{Streamer: ap.resampler, Base: 2}
	ap.mediaURI = uri
	ap.index = index
	ap.total = total
	ap.artist = artist
	ap.title = title
	ap.done = done
}

func (ap *AudioPanel) SetMessage(message string) {
	ap.message = message
}

func (ap *AudioPanel) IsPaused() bool {
	return ap.ctrl.Paused
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

	s := fmt.Sprintf("Media   [%d/%d] (P/N):", ap.index, ap.total)
	drawTextLine(screen, 0, 5, s, mainStyle)
	drawTextLine(screen, len(s), 5, ap.mediaURI, statusStyle)

	row := 6
	if ap.title != "" {
		drawTextLine(screen, 0, row, "Title"+strings.Repeat(" ", len(s)-len(`Title:`))+":", mainStyle)
		drawTextLine(screen, len(s), row, ap.title, statusStyle)
		row++
	}
	if ap.artist != "" {
		drawTextLine(screen, 0, row, "Artist"+strings.Repeat(" ", len(s)-len(`Artist:`))+":", mainStyle)
		drawTextLine(screen, len(s), row, ap.artist, statusStyle)
		row++
	}

	drawTextLine(screen, 0, row, "Position"+strings.Repeat(" ", len(s)-len(`Position(Q/W):`))+"(Q/W):", mainStyle)
	drawTextLine(screen, len(s), row, positionStatus, statusStyle)
	row++

	drawTextLine(screen, 0, row, "Volume"+strings.Repeat(" ", len(s)-len(`Volume(A/S):`))+"(A/S):", mainStyle)
	drawTextLine(screen, len(s), row, volumeStatus, statusStyle)
	row++

	drawTextLine(screen, 0, row, "Speed"+strings.Repeat(" ", len(s)-len(`Speed(Z/X):`))+"(Z/X):", mainStyle)
	drawTextLine(screen, len(s), row, speedStatus, statusStyle)
	row++

	drawTextLine(screen, 0, row, "Repeat/Shuffle"+strings.Repeat(" ", len(s)-len(`Repeat/Shuffle(R/F):`))+"(R/F):", mainStyle)
	drawTextLine(screen, len(s), row, fmt.Sprintf("%s/%s", util.Bool2Str(config.Repeat), util.Bool2Str(config.Shuffle)), statusStyle)
	row++

	if ap.message != "" {
		drawTextLine(screen, 0, row+1, ap.message, mainStyle)
	}
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
			if ap.streamer.Len() > 0 {
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

		case 'r':
			config.Repeat = !config.Repeat
			return true, HandleActionNOP

		case 'f':
			config.Shuffle = !config.Shuffle
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
	case *tcell.EventResize:
		return true, HandleActionNOP
	}
	return false, HandleActionNOP
}
