package id_generator

import (
	"encoding/binary"
	"strconv"

	"github.com/google/uuid"

	"github.com/zitadel/zitadel/internal/id_generator/sonyflake"
	"github.com/zitadel/zitadel/internal/id_generator/uuidv7"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GeneratorType int

const (
	GeneratorTypeMock GeneratorType = -1

	GeneratorTypeSonyFlake GeneratorType = 0
	GeneratorTypeUUIDv7    GeneratorType = 1
)

func (g GeneratorType) IsValid() bool {
	switch g {
	case GeneratorTypeMock, GeneratorTypeSonyFlake, GeneratorTypeUUIDv7:
		return true
	default:
		return false
	}
}

type Generator interface {
	Next() (string, error)
}

var (
	generator Generator = nil
)

func Next() (string, error) {
	if generator == nil {
		panic("generator not configured")
	}

	return generator.Next()
}

func SetGenerator(g Generator) {
	generator = g
}

func SetGeneratorWithConfig(t GeneratorType, c ...interface{}) {
	if !t.IsValid() {
		panic("invalid generator type")
	}

	switch t {
	case GeneratorTypeSonyFlake:
		if len(c) != 1 {
			panic("invalid SonyFlake config")
		}
		config, ok := c[0].(*sonyflake.Config)
		if !ok {
			panic("invalid SonyFlake config type")
		}

		SetGenerator(sonyflake.New(config))
	case GeneratorTypeUUIDv7:
		SetGenerator(&uuidv7.Generator{})
	default:
		panic("unsupported generator type")
	}
}

func NumericFromID(id string) (int64, error) {
	// if the id is a number, we can use it directly
	if n, err := strconv.Atoi(id); err == nil {
		return int64(n), nil
	}

	// if the id is a UUIDv7, we can use the first 8 bytes as numeric id, as they are monotonically increasing
	if parse, err := uuid.Parse(id); err == nil && parse.Version() == 7 {
		return int64(binary.BigEndian.Uint64(parse[:8])), nil
	}

	return 0, zerrors.ThrowInternalf(nil, "ID-40eedf83", "failed to get numeric id from id %s", id)
}
