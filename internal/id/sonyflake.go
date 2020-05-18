package id

import (
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

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

var (
	SonyFlakeGenerator = Generator(&sonyflakeGenerator{
		sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: time.Date(2019, 4, 29, 0, 0, 0, 0, time.UTC),
		}),
	})
)
