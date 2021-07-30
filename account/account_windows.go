package account

import (
	"github.com/StackExchange/wmi"
	"time"
)

type Group struct {
	Caption      string    `json:"caption"`
	Description  string    `json:"description"`
	Domain       string    `json:"domain"`
	InstallDate  time.Time `json:"install_date"`
	LocalAccount bool      `json:"local_account"`
	Name         string    `json:"name"`
	Sid          string    `json:"sid"`
	SidType      uint8     `json:"sid_type"`
	Status       string    `json:"status"`
}

type Account struct {
	AccountType        uint32    `json:"account_type"`
	Caption            string    `json:"caption"`
	Description        string    `json:"description"`
	Disabled           bool      `json:"disabled"`
	Domain             string    `json:"domain"`
	FullName           string    `json:"full_name"`
	InstallDate        time.Time `json:"install_date"`
	LocalAccount       bool      `json:"local_account"`
	Lockout            bool      `json:"lockout"`
	Name               string    `json:"name"`
	PasswordChangeable bool      `json:"password_changeable"`
	PasswordExpires    bool      `json:"password_expires"`
	PasswordRequired   bool      `json:"password_required"`
	SID                string    `json:"sid"`
	SIDType            uint8     `json:"sid_type"`
	Status             string    `json:"status"`
}

type Accounts []*Account
type Groups []*Group

var (
	WQLAccount = "SELECT * FROM Win32_UserAccount"
	WQLGroup   = "SELECT * FROM Win32_Account"
)

// GetAll 获取用户和用户组信息
func GetAll() *Info {
	accounts := GetAccounts()
	groups := GetGroups()
	return &Info{
		Accounts: accounts,
		Groups:   groups,
	}
}

func GetAccounts() Accounts {
	var dst Accounts
	err := wmi.Query(WQLAccount, &dst)
	if err != nil {
		return nil
	}

	return dst
}

func GetGroups() Groups {
	var dst Groups
	err := wmi.Query(WQLGroup, &dst)
	if err != nil {
		return nil
	}

	return dst
}
