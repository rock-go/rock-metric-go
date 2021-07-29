package command

import "github.com/rock-go/rock/logger"

func GetHistory(user string) map[string][]*History {
	logger.Errorf("not supported on Windows")
	return nil
}
