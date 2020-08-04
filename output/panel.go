package output

import (
	"fmt"
	"strings"
	"sync/atomic"
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
	HandleActionDownload
	HandleActionM3U
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
	screen   tcell.Screen
	mediaURI string
	index    int
	total    int
	artist   string
	title    string
	message  string
	quit     atomic.Value
}

func NewScreenPanel() *ScreenPanel {
	sp := &ScreenPanel{}
	sp.quit.Store(false)
	return sp
}

func (sp *ScreenPanel) PollScreenEvent() tcell.Event {
	return sp.screen.PollEvent()
}

func (sp *ScreenPanel) Quit() bool {
	return sp.quit.Load().(bool)
}

func (sp *ScreenPanel) Initialize() error {
	if sp.screen != nil {
		return nil
	}
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	err = screen.Init()
	if err != nil {
		return err
	}

	sp.screen = screen
	return nil
}

func (sp *ScreenPanel) Finalize() {
	sp.quit.Store(true)
	if sp.screen != nil {
		sp.screen.Fini()
		sp.screen = nil
	}
}

func (sp *ScreenPanel) Update(uri string, index int, total int, artist string, title string) {
	sp.mediaURI = uri
	sp.index = index
	sp.total = total
	sp.artist = artist
	sp.title = title
}

func (sp *ScreenPanel) SetMessage(message string) {
	sp.message = message
}

func (sp *ScreenPanel) Draw(position time.Duration, length time.Duration, volume float64, speed float64) {
	sp.screen.Clear()
	mainStyle := tcell.StyleDefault.
		Background(tcell.NewHexColor(0x473437)).
		Foreground(tcell.NewHexColor(0xD7D8A2))
	statusStyle := mainStyle.
		Foreground(tcell.NewHexColor(0xDDC074)).
		Bold(true)

	sp.screen.Fill(' ', mainStyle)

	drawTextLine(sp.screen, 0, 0, "Welcome to Hannah, the simplest music player!", mainStyle)
	drawTextLine(sp.screen, 0, 1, "Press [ESC] to quit.", mainStyle)
	drawTextLine(sp.screen, 0, 2, "Press [SPACE] to pause/resume.", mainStyle)
	drawTextLine(sp.screen, 0, 3, "Use keys in (?/?) to turn the buttons.", mainStyle)

	positionStatus := fmt.Sprintf("%v / %v", position.Round(time.Second), length.Round(time.Second))
	volumeStatus := fmt.Sprintf("%.1f", volume)
	speedStatus := fmt.Sprintf("%.3fx", speed)

	s := fmt.Sprintf("Media   [%d/%d] (P/N):", sp.index, sp.total)
	drawTextLine(sp.screen, 0, 5, s, mainStyle)
	drawTextLine(sp.screen, len(s), 5, sp.mediaURI, statusStyle)

	row := 6
	if sp.title != "" {
		drawTextLine(sp.screen, 0, row, "Title"+strings.Repeat(" ", len(s)-len(`Title:`))+":", mainStyle)
		drawTextLine(sp.screen, len(s), row, sp.title, statusStyle)
		row++
	}
	if sp.artist != "" {
		drawTextLine(sp.screen, 0, row, "Artist"+strings.Repeat(" ", len(s)-len(`Artist:`))+":", mainStyle)
		drawTextLine(sp.screen, len(s), row, sp.artist, statusStyle)
		row++
	}

	drawTextLine(sp.screen, 0, row, "Position"+strings.Repeat(" ", len(s)-len(`Position(Q/W):`))+"(Q/W):", mainStyle)
	drawTextLine(sp.screen, len(s), row, positionStatus, statusStyle)
	row++

	drawTextLine(sp.screen, 0, row, "Volume"+strings.Repeat(" ", len(s)-len(`Volume(A/S):`))+"(A/S):", mainStyle)
	drawTextLine(sp.screen, len(s), row, volumeStatus, statusStyle)
	row++

	drawTextLine(sp.screen, 0, row, "Speed"+strings.Repeat(" ", len(s)-len(`Speed(Z/X):`))+"(Z/X):", mainStyle)
	drawTextLine(sp.screen, len(s), row, speedStatus, statusStyle)
	row++

	drawTextLine(sp.screen, 0, row, "Repeat/Shuffle"+strings.Repeat(" ", len(s)-len(`Repeat/Shuffle(R/F):`))+"(R/F):", mainStyle)
	drawTextLine(sp.screen, len(s), row, fmt.Sprintf("%s/%s", util.Bool2Str(config.Repeat), util.Bool2Str(config.Shuffle)), statusStyle)
	row++

	drawTextLine(sp.screen, 0, row, "Download/M3U"+strings.Repeat(" ", len(s)-len(`Download/M3U(D/V):`))+"(D/V):", mainStyle)
	drawTextLine(sp.screen, len(s), row, fmt.Sprintf("Download current song to %s/Add current song to %s", config.DownloadDir, config.M3UFileName), statusStyle)
	row++

	if sp.message != "" {
		drawTextLine(sp.screen, 0, row+1, sp.message, mainStyle)
	}

	sp.screen.Show()
}

func (sp *ScreenPanel) Handle(event tcell.Event) (changed bool, action int) {
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
			'd': HandleActionDownload,
			'v': HandleActionM3U,
		}
		if action, ok := cmdMap[unicode.ToLower(event.Rune())]; ok {
			return true, action
		}
	case *tcell.EventResize:
		return true, HandleActionNOP
	}
	return false, HandleActionNOP
}
