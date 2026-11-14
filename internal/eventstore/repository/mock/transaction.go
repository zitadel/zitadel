package mock

import "context"

// Unit command fixtures run synchronously; database tests cover transaction locks.
func (m *MockPusher) WithInstanceLock(ctx context.Context, change func(context.Context) error) error {
	return change(ctx)
}
