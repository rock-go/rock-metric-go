package diskIo

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDetail(t *testing.T) {
	for {
		s := GetDetail()
		if len(s) < 2 {
			continue
		}
		//fmt.Printf("%v\n", s)
		//rs := s[0].ReadPerSecBytes + s[1].ReadPerSecBytes
		ws := s[0].WritePerSecByte + s[1].WritePerSecByte
		fmt.Printf("%s,%f,\n", time.Now().Format("04:05"), ws)
		time.Sleep(1 * time.Second)
	}
}
