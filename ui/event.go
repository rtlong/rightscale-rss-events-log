package ui

type Event struct {
	Type int
	Err  error
}

const (
	EventQuit int = iota
	EventError
)
