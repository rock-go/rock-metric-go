package command

import (
	"fmt"
	"testing"
)

func TestGetHistory(t *testing.T) {
	//GetByCMD("root")
	//h := GetFromFile("root")
	d := GetHistory("")
	fmt.Println(d)
}
