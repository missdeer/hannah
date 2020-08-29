package lyric

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func LRC2SMI(src string) string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(src))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		ll := r3.FindAllStringSubmatch(line, -1)
		if len(ll) > 0 && len(ll[0]) == 5 {
			m, err := strconv.Atoi(ll[0][1])
			if err != nil {
				continue
			}
			s, err := strconv.Atoi(ll[0][2])
			if err != nil {
				continue
			}
			ms, err := strconv.Atoi(ll[0][3])
			if err != nil {
				continue
			}
			t := ms + s*1000 + m*60*1000
			l := fmt.Sprintf(lineTemplate, t, ll[0][4])
			lines = append(lines, l)
		}
	}
	return fmt.Sprintf(t, strings.Join(lines, "\n"))
}
