package uuid

import "github.com/google/uuid"

type uuidGenerator struct {
}

func NewUUIDGenerator() *uuidGenerator {

	return &uuidGenerator{}
}

func (u *uuidGenerator) GenUID() string {
	// uuid V4
	return uuid.NewString()
}
