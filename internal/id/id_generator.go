package id

import (
	"strconv"

	"github.com/sony/sonyflake"
)

type Generator interface {
	Next() (string, error)
}

type sonyflakeGenerator struct {
	*sonyflake.Sonyflake
}

func (s *sonyflakeGenerator) Next() (string, error) {
	id, err := s.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

func NewSonyFlake(s *sonyflake.Sonyflake) Generator {
	return &sonyflakeGenerator{
		s,
	}
}
