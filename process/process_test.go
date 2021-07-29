package process

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGetProc(t *testing.T) {
	proc, _ := newProcess(14768)
	fmt.Println(proc)
	m := Metric{}
	fmt.Println(m)
}

func TestGetSummary(t *testing.T) {
	for {
		s := GetSummary("")
		data, _ := json.Marshal(s)
		fmt.Printf("%s\n", data)
		time.Sleep(1 * time.Second)
	}
}
