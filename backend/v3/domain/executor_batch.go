package domain

import (
	"context"
	"fmt"
	"strings"
)

// batchExecutor is a batch of [Executor]s.
// It uses the [Invoker]s provided by the opts to execute each [Executor].
// The [Executor] is sent to all [Invoker]s before moving to the next [Executor].
type batchExecutor struct {
	executors []Executor
}

func BatchExecutor(executors ...Executor) *batchExecutor {
	return &batchExecutor{
		executors: executors,
	}
}

// String implements [Executor].
func (cmd *batchExecutor) String() string {
	names := make([]string, len(cmd.executors))
	for i, c := range cmd.executors {
		names[i] = c.String()
	}
	return fmt.Sprintf("commandBatch[%s]", strings.Join(names, ", "))
}

// Execute implements [Executor].
func (b *batchExecutor) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	for _, cmd := range b.executors {
		if err = opts.Invoke(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

var _ Executor = (*batchExecutor)(nil)
