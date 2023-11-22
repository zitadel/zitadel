package mock

import (
	"context"
	"io"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/static"
)

func NewStorage(t *testing.T) *MockStorage {
	return NewMockStorage(gomock.NewController(t))
}

func (m *MockStorage) ExpectPutObject() *MockStorage {
	m.EXPECT().
		PutObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, instanceID, location, resourceOwner, name, contentType string, objectType static.ObjectType, object io.Reader, objectSize int64) (*static.Asset, error) {
			hash, _ := io.ReadAll(object)
			return &static.Asset{
				InstanceID:   instanceID,
				Name:         name,
				Hash:         string(hash),
				Size:         objectSize,
				LastModified: time.Now(),
				Location:     location,
				ContentType:  contentType,
			}, nil
		})
	return m
}

func (m *MockStorage) ExpectPutObjectError() *MockStorage {
	m.EXPECT().
		PutObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, caos_errors.ThrowInternal(nil, "", ""))
	return m
}

func (m *MockStorage) ExpectRemoveObjectNoError() *MockStorage {
	m.EXPECT().
		RemoveObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
		RemoveObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(caos_errors.ThrowInternal(nil, "", ""))
	return m
}
