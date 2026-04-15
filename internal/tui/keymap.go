package tui

type KeyMap struct {
	Up    string
	Down  string
	Open  string
	Close string
	Quit  string
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:    "↑/k",
		Down:  "↓/j",
		Open:  "Enter",
		Close: "Esc",
		Quit:  "Ctrl+C",
	}
}
