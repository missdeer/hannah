package action

import (
	"github.com/missdeer/hannah/config"
)

type actionHandler func(...string) error

var (
	actionHandlerMap = map[string]actionHandler{
		"play":     play,
		"search":   search,
		"m3u":      save,
		"download": download,
		"hot":      hot,
		"playlist": playlist,
	}
)

func GetActionHandler(action string) actionHandler {
	return actionHandlerMap[config.Action]
}
