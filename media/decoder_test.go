package media

import (
	"testing"
)

func TestBuiltinSupportedFileType(t *testing.T) {
	if !BuiltinSupportedFileType(".mp3") {
		t.Error("mp3 should be supported")
	}
	if !BuiltinSupportedFileType(".flac") {
		t.Error("flac should be supported")
	}
	if !BuiltinSupportedFileType(".ogg") {
		t.Error("ogg should be supported")
	}
	if !BuiltinSupportedFileType(".wav") {
		t.Error("wav should be supported")
	}

	if BuiltinSupportedFileType(".mp4") {
		t.Error("mp4 should not be supported")
	}
	if BuiltinSupportedFileType(".m4a") {
		t.Error("m4a should not be supported")
	}
}
