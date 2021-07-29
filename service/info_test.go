package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestGetWinService(t *testing.T) {
	s := GetService("")
	data, _ := json.Marshal(s)

	fmt.Printf("%s\n", data)
}

func TestStr(t *testing.T) {
	a := " "
	fmt.Println(strings.Contains("aaabbb", a))
}
