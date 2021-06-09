package mock

import (
	"testing"

	"github.com/golang/mock/gomock"

	caos_errors "github.com/caos/zitadel/internal/errors"
)

func NewStorage(t *testing.T) *MockStorage {
	return NewMockStorage(gomock.NewController(t))
}

func (m *MockStorage) ExpectAddObjectNoError() *MockStorage {
	m.EXPECT().
		PutObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)
	return m
}

func (m *MockStorage) ExpectAddObjectError() *MockStorage {
	m.EXPECT().
		PutObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, caos_errors.ThrowInternal(nil, "", ""))
	return m
}

func (m *MockStorage) ExpectRemoveObjectNoError() *MockStorage {
	m.EXPECT().
		RemoveObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	return m
}

func (m *MockStorage) ExpectRemoveObjectsNoError() *MockStorage {
	m.EXPECT().
		RemoveObjects(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	return m
}

func (m *MockStorage) ExpectRemoveObjectError() *MockStorage {
	m.EXPECT().
		RemoveObject(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(caos_errors.ThrowInternal(nil, "", ""))
	return m
}
