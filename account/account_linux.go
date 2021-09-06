package account

import (
	"bufio"
	"github.com/rock-go/rock/logger"
	"io"
	"os"
	"strings"
)

type Account struct {
	LoginName string `json:"login_name"`
	UID       string `json:"uid"`
	GID       string `json:"gid"`
	UserName  string `json:"user_name"`
	HomeDir   string `json:"home_dir"`
	Shell     string `json:"shell"`
	Raw       string `json:"raw"`
}

type Group struct {
	GroupName string `json:"group_name"`
	GID       string `json:"gid"`
	Raw       string `json:"raw"`
}

type Accounts []*Account
type Groups []*Group

var (
	ACCOUNT = "/etc/passwd"
	GROUP   = "/etc/group"
)

func GetAll() *Info {
	var info Info
	info.Accounts = GetAccounts()
	info.Groups = GetGroups()
	return &info
}

func GetAccounts() Accounts {
	f, err := os.OpenFile(ACCOUNT, os.O_RDONLY, 0666)
	if err != nil {
		logger.Errorf("read /etc/passwd error: %v", err)
		return nil
	}

	var accounts Accounts
	rd := bufio.NewReaderSize(f, 4096)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}

		account := &Account{Raw: strings.Trim(line, "\n")}
		data := strings.Split(line, ":")
		if len(data) < 7 {
			goto APPEND
		}
		account.LoginName = data[0]
		account.UID = data[2]
		account.GID = data[3]
		account.UserName = data[4]
		account.HomeDir = data[5]
		account.Shell = strings.Trim(data[6], "\n")
	APPEND:
		accounts = append(accounts, account)
	}

	return accounts
}

func GetGroups() Groups {
	f, err := os.OpenFile(GROUP, os.O_RDONLY, 0666)
	if err != nil {
		logger.Errorf("read /etc/group error: %v", err)
		return nil
	}

	var groups Groups
	rd := bufio.NewReaderSize(f, 4096)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}

		group := &Group{Raw: strings.Trim(line, "\n")}
		data := strings.Split(line, ":")
		if len(data) < 4 {
			goto APPEND
		}
		group.GroupName = data[0]
		group.GID = data[2]
	APPEND:
		groups = append(groups, group)
	}

	return groups
}
