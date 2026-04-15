package tui

type KeyMap struct {
	Up     string
	Down   string
	Search string
	Open   string
	Close  string
	Quit   string
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:     "↑/k",
		Down:   "↓/j",
		Search: "Enter",
		Open:   "Ctrl+O",
		Close:  "Esc",
		Quit:   "Ctrl+C",
	}
}
