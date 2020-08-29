package lyric

func LyricConvert(fromFormat string, toFormat string, src string) (dst string) {
	switch fromFormat {
	case "lrc":
		switch toFormat {
		case "smi":
			return LRC2SMI(src)
		default:
			return src
		}
	case "xtrc":
		switch toFormat {
		case "lrc":
			return XTRC2LRC(src)
		case "smi":
			return XTRC2SMI(src)
		default:
			return src
		}
	default:
		return src
	}
}
