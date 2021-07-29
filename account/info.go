package account

type Info struct {
	Accounts []*Account
	Groups   []*Group
}

func GetInfo() *Info {
	return GetAll()
}
