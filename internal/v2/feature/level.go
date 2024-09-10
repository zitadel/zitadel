package feature

import (
	feature "github.com/zitadel/zitadel/internal/feature"
)

type Level = feature.Level

var Levels = feature.LevelValues()

func levelFromString(text string) Level {
	level, _ := feature.LevelString(text)
	return level
}
