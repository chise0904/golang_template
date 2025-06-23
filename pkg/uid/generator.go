package uid

import (
	"github.com/chise0904/golang_template/pkg/uid/intid"
	"github.com/chise0904/golang_template/pkg/uid/mongoid"
	"github.com/chise0904/golang_template/pkg/uid/uuid"
	"github.com/chise0904/golang_template/pkg/uid/xid"
)

type generatorEnum int8

const (
	GeneratorEnumXID generatorEnum = 0 + iota //default
	GeneratorEnumUUID
	GeneratorEnumMongoID
	GeneratorEnumIntID
)

type UIDGenerator interface {
	GenUID() string
}

func NewUIDGenerator(e generatorEnum) UIDGenerator {
	switch e {
	case GeneratorEnumXID:
		return xid.NewXIDGenerator()
	case GeneratorEnumUUID:
		return uuid.NewUUIDGenerator()
	case GeneratorEnumMongoID:
		return mongoid.NewMongoIDGenerator()
	case GeneratorEnumIntID:
		return intid.NewIntIDGenerator()
	default:
		return xid.NewXIDGenerator()
	}
}
