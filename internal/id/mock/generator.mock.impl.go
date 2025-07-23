package mock

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func NewIDGenerator(t *testing.T) *MockGenerator {
	t.Helper()

	m := NewMockGenerator(gomock.NewController(t))
	m.EXPECT().Next().Return("1", nil)
	return m
}

func NewIDGeneratorExpectIDs(t *testing.T, ids ...string) *MockGenerator {
	t.Helper()

	m := NewMockGenerator(gomock.NewController(t))
	for _, id := range ids {
		m.EXPECT().Next().Return(id, nil)
	}
	return m
}

func ExpectID(t *testing.T, id string) *MockGenerator {
	t.Helper()

	m := NewMockGenerator(gomock.NewController(t))
	m.EXPECT().Next().Return(id, nil)
	return m
}

func NewIDGeneratorExpectError(t *testing.T, err error) *MockGenerator {
	t.Helper()

	m := NewMockGenerator(gomock.NewController(t))
	m.EXPECT().Next().Return("", err)
	return m
}
