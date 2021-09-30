package network

import (
	"github.com/rock-go/rock/lua"
)

type detail []Ifc

func (d *detail) Byte() []byte {
	return Json(*d)
}

func (d *detail) String() string {
	return lua.B2S(d.Byte())
}