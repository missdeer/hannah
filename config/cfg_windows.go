// +build windows

package config

var (
	AudioDriver     = "WinMM"
	AudioDriverList = []string{"ASIO", "WASAPI", "DirectSound", "WinMM"}
)
