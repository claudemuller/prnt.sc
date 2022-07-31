package pkg

import (
	"gioui.org/app"
)

type State struct {
	Win        *app.Window
	IDs        []string
	MaxRetries *int
}

func NewState(maxRetries *int) *State {
	return &State{
		Win:        nil,
		IDs:        make([]string, 1),
		MaxRetries: maxRetries,
	}
}
