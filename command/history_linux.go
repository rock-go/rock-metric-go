//+build linux
// linux history commands

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/rock-go/rock-metric-go/account"
	"github.com/rock-go/rock/logger"
	"io"
	"os"
	"os/exec"
	"strings"
)

func GetHistory(userName string) HistoryMap {
	hm := make(HistoryMap)
	users := account.GetAccounts()
	for _, u := range users {
		if userName == "" || u.UserName == userName {
			hm[u.UserName] = GetFromFile(*u)
		}
	}

	return hm
}

func GetFromFile(user account.Account) []*History {
	homeDir := user.HomeDir
	if homeDir == "" {
		return nil
	}

	f, err := os.OpenFile(homeDir+"/.bash_history", os.O_RDONLY, 0666)
	if err != nil {
		logger.Errorf("read .bash_history error: %v", err)
		return nil
	}

	var histories []*History
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}

		history := &History{
			User:    user.UserName,
			ID:      "",
			Command: strings.Trim(line, "\n"),
		}

		histories = append(histories, history)
	}

	return histories
}

func GetByCMD(user string) []*History {
	var history = make([]*History, 0)
	var buff bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("su", user)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		logger.Errorf("execute command [%s] error: %v", cmd.String(), err)
		logger.Errorf("stderr: %s", stderr.String())
		return nil
	}

	cmd = exec.Command("", "-c", "history")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		logger.Errorf("execute command [%s] error: %v", cmd.String(), err)
		logger.Errorf("stderr: %s", stderr.String())
		return nil
	}
	buff.Write(stdout.Bytes())

	fmt.Println(buff.String())
	return history
}
