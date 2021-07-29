package service

type Service struct {
	Name        string `json:"name"`
	StartType   string `json:"start_type"`
	ExecPath    string `json:"exec_path"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	State       string `json:"state"`
	Pid         uint32 `json:"pid"`
	ExitCode    uint32 `json:"exit_code"`
}

type Services []*Service

func (ss *Services) DisableReflect() {}

func GetDetail(pattern string) Services {
	return Services(GetService(pattern))
}
