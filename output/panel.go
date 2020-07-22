package output

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gdamore/tcell"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/util"
)

const (
	HandleActionQUIT = iota
	HandleActionPREVIOUS
	HandleActionNEXT
	HandleActionNOP
	HandleActionPauseResume
	HandleActionBackward
	HandleActionForward
	HandleActionSlowdown
	HandleActionSpeedup
	HandleActionIncreaseVolume
	HandleActionDecreaseVolume
	HandleActionRepeat
	HandleActionShuffle
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

type ScreenPanel struct {
	mediaURI string
	index    int
	total    int
	artist   string
	title    string
	message  string
}

func NewScreenPanel(uri string, index int, total int, artist string, title string) *ScreenPanel {
	return &ScreenPanel{
		mediaURI: uri,
		index:    index,
		total:    total,
		artist:   artist,
		title:    title,
	}
}

func (ap *ScreenPanel) Update(uri string, index int, total int, artist string, title string) {
	ap.mediaURI = uri
	ap.index = index
	ap.total = total
	ap.artist = artist
	ap.title = title
}

func (ap *ScreenPanel) SetMessage(message string) {
	ap.message = message
}

func (ap *ScreenPanel) Draw(screen tcell.Screen, position time.Duration, length time.Duration, volume float64, speed float64) {
	screen.Clear()
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

	screen.Show()
}

func (ap *ScreenPanel) Handle(event tcell.Event) (changed bool, action int) {
	switch event := event.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyESC {
			return false, HandleActionQUIT
		}

		if event.Key() != tcell.KeyRune {
			return false, HandleActionNOP
		}

		cmdMap := map[rune]int{
			' ': HandleActionPauseResume,
			'q': HandleActionBackward,
			'w': HandleActionForward,
			'a': HandleActionDecreaseVolume,
			's': HandleActionIncreaseVolume,
			'z': HandleActionSlowdown,
			'x': HandleActionSpeedup,
			'p': HandleActionPREVIOUS,
			'n': HandleActionNEXT,
			'r': HandleActionRepeat,
			'f': HandleActionShuffle,
		}
		if action, ok := cmdMap[unicode.ToLower(event.Rune())]; ok {
			return true, action
		}
	case *tcell.EventResize:
		return true, HandleActionNOP
	}
	return false, HandleActionNOP
}
