package xid

import (
	"github.com/rs/xid"
)

type xidGenerator struct {
	x xid.ID
}

func NewXIDGenerator() *xidGenerator {
	return &xidGenerator{
		x: xid.New(),
	}
}

func (x *xidGenerator) GenUID() string {
	return x.x.String()
}
