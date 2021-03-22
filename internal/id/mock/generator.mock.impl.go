package mock

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func NewIDGenerator(t *testing.T) *MockGenerator {
	m := NewMockGenerator(gomock.NewController(t))
	m.EXPECT().Next().Return("1", nil)
	return m
}

func NewIDGeneratorExpectIDs(t *testing.T, ids ...string) *MockGenerator {
	m := NewMockGenerator(gomock.NewController(t))
	for _, id := range ids {
		m.EXPECT().Next().Return(id, nil)
	}
	return m
}

func ExpectID(t *testing.T, id string) *MockGenerator {
	m := NewMockGenerator(gomock.NewController(t))
	m.EXPECT().Next().Return(id, nil)
	return m
}
