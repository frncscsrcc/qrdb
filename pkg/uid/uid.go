package uid

import (
	"qrdb/pkg/libs/random"
)

type UIDGenerator interface {
	GetUID() string
}

type BasicUIDGenerator struct {
	len int
}

func NewBasicUIDGenerator(len int) BasicUIDGenerator {
	return BasicUIDGenerator{len}
}

func (ug BasicUIDGenerator) GetUID() string {
	return random.GetRandomString(ug.len)
}
