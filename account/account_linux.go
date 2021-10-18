package account

import (
	"bufio"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
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

func GetAccounts() As {
	f, err := os.OpenFile(ACCOUNT, os.O_RDONLY, 0666)
	if err != nil {
		logger.Errorf("read /etc/passwd error: %v", err)
		return nil
	}
	defer f.Close()

	var accounts As
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

func GetGroups() Gs {
	f, err := os.OpenFile(GROUP, os.O_RDONLY, 0666)
	if err != nil {
		logger.Errorf("read /etc/group error: %v", err)
		return nil
	}
	defer f.Close()

	var groups Gs
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

func (a As) Byte() []byte {
	buf := json.NewBuffer()
	buf.Arr("")
	for _, item := range a {
		buf.Tab("")
		buf.KV("login_name", item.LoginName)
		buf.KV("uid", item.UID)
		buf.KV("gid", item.GID)
		buf.KV("user_name", item.UserName)
		buf.KV("home_dir", item.HomeDir)
		buf.KV("shell", item.Shell)
		buf.KV("raw", item.Raw)
		buf.End("},")
	}

	buf.End("]")

	return buf.Bytes()
}

func (a As) String() string {
	return lua.B2S(a.Byte())
}

func (g Gs) Byte() []byte {
	buf := json.NewBuffer()
	for _, item := range g {
		buf.WriteByte('{')
		buf.KV("group_name", item.GroupName)
		buf.KV("gid", item.GID)
		buf.KV("raw", item.Raw)
		buf.WriteString("},")
	}
	buf.End("]")
	return buf.Bytes()
}

func (g Gs) String() string {
	return lua.B2S(g.Byte())
}
