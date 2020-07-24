package action

import (
	"testing"
)

func TestGetActionHandler(t *testing.T) {
	h := GetActionHandler("play")
	if h == nil {
		t.Error("can't get play action handler")
	}
	h = GetActionHandler("search")
	if h == nil {
		t.Error("can't get search action handler")
	}
}
