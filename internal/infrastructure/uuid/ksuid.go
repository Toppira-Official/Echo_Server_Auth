package uuid

import "github.com/segmentio/ksuid"

type KsuidIdGenerator struct{}

func NewKsuidIdGenerator() *KsuidIdGenerator { return &KsuidIdGenerator{} }

func (*KsuidIdGenerator) Generate() (string, error) {
	return ksuid.New().String(), nil
}
