package socket

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGetSockets(t *testing.T) {
	s := Summary{}
	GetSockets(&s, "")
	fmt.Println(s)
}

func TestName(t *testing.T) {

	n := Summary{}
	// get
	immutable := reflect.ValueOf(n)
	val := immutable.FieldByName("LISTEN").Int()
	fmt.Printf("N=%d\n", val) // prints 1

	// set
	mutable := reflect.ValueOf(&n).Elem()
	mutable.FieldByName("LISTEN").SetInt(7)
	fmt.Printf("N=%d\n", n.LISTEN) // prints 7
}

func TestMarshal(t *testing.T) {
	s1 := &Socket{
		State:      "",
		LocalIP:    "",
		LocalPort:  0,
		RemoteIP:   "",
		RemotePort: 0,
		Pid:        0,
	}
	s2 := &Socket{
		State:      "",
		LocalIP:    "",
		LocalPort:  0,
		RemoteIP:   "",
		RemotePort: 0,
		Pid:        0,
	}
	s := Summary{
		CLOSED:      0,
		LISTEN:      0,
		SYN_SENT:    0,
		SYN_RCVD:    0,
		ESTABLISHED: 0,
		FIN_WAIT1:   0,
		FIN_WAIT2:   0,
		CLOSE_WAIT:  0,
		CLOSING:     0,
		LAST_ACK:    0,
		TIME_WAIT:   0,
		DELETE_TCB:  0,
		Sockets:     []*Socket{s1, nil, s2},
	}

	data, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
