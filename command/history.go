package command

type History struct {
	User    string `json:"user"`
	ID      string `json:"id"`
	Command string `json:"command"`
}

type HistoryMap map[string][]*History