package domain

import (
	"context"
	"fmt"
	"strings"
)

// executorBatch is a batch of commands.
// It uses the [Invoker] provided by the opts to execute each command.
type executorBatch struct {
	executors []Executor
}

func BatchExecutor(executors ...Executor) *executorBatch {
	return &executorBatch{
		executors: executors,
	}
}

// String implements [Executor].
func (cmd *executorBatch) String() string {
	names := make([]string, len(cmd.executors))
	for i, c := range cmd.executors {
		names[i] = c.String()
	}
	return fmt.Sprintf("commandBatch[%s]", strings.Join(names, ", "))
}

// Execute implements [Executor].
func (b *executorBatch) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	for _, cmd := range b.executors {
		if err = opts.Invoke(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

var _ Executor = (*executorBatch)(nil)
