package action

import (
	"github.com/missdeer/hannah/config"
)

type actionHandler func([]string) error

var (
	actionHandlerMap = map[string]actionHandler{
		"play":   play,
		"search": search,
	}
)

func GetActionHandler(action string) actionHandler {
	return actionHandlerMap[config.Action]
}
