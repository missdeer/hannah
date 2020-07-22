// +build freebsd netbsd openbsd

package config

var (
	AudioDriver     = "OSS"
	AudioDriverList = []string{"OSS", "audio", "sndio"}
)
