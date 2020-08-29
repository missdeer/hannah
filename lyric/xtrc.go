package lyric

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	r1 = regexp.MustCompile(`<\d+>`)
	r2 = regexp.MustCompile(`\[x\-trans\][^\n]*\n`)
	r3 = regexp.MustCompile(`^\[(\d{2}):(\d{2})\.(\d{2,3})\]([^\n]+)`)

	t = `
<SAMI>
<head>
<style type="text/css"><!--
p     { background-color: #000000; text-align: center; font-size: 16pt; font-family: SimSun; }
.QCRJ { name: Chinese; lang: ZH-CN;  sami_type: CC; }
--></style>
</head>
<body>
<sync start=0><p class=QCRJ><font color=#FFFFFF> 
%s
<sync start=232672000><p class=QCRJ><font color=#FFFFFF> 
</body>
</SAMI>
`
	lineTemplate = `<sync start=%d><p class=QCRJ><font color=#FFFFFF>%s`
)

func XTRC2LRC(src string) (dst string) {
	dst = r1.ReplaceAllString(src, "")
	dst = r2.ReplaceAllString(dst, "")
	return
}

func XTRC2SMI(src string) (dst string) {
	dst = r1.ReplaceAllString(src, "")
	dst = r2.ReplaceAllString(dst, "")

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(dst))
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
