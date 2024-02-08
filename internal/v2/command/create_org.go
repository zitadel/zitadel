package command

import "github.com/zitadel/zitadel/internal/eventstore"

/*
api/org/v2 -> internal/v2/org/Create

*/

type Create struct {
	eventstore.WriteModel

	Name   string
	Domain string
}

type previousSequence func(currentSequence uint64) bool

func DoesntMatter(uint64) bool            { return true }
func AtLeast(currentSequence uint64) bool { return currentSequence >= 10 }
func Exactly(currentSequence uint64) bool { return currentSequence == 10 }

func (c *Create) PreviousSequence() previousSequence {
	return func(currentSequence uint64) bool {
		return currentSequence == 0
	}
}
