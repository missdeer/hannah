package util

func Bool2Str(b bool) string {
	if b {
		return "Enabled"
	}
	return "Disabled"
}
