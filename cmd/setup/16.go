package setup

import (
	"context"
)

type RepeatStep10 struct {
	step10 *CorrectCreationDate
}

func (mig *RepeatStep10) Execute(ctx context.Context) error {
	// execute step 10 again because events created after the first execution of step 10
	// could still have the wrong ordering of sequences and creation date
	return mig.step10.Execute(ctx)
}

func (mig *RepeatStep10) String() string {
	return "16_repeat_correct_creation_date"
}
