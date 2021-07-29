package cpu

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	//i := 0
	for i := 0; i < 2; i++ {
		go func() {
			for {

			}
		}()
	}

	M := Metric{}
	for {
		cpu := Get(&M)
		if cpu != nil {
			fmt.Println(cpu.Total, cpu.Idle)
		}
		time.Sleep(1 * time.Second)
	}
}
