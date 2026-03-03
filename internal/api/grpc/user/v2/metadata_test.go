package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func Test_setUserMetadataToDomain(t *testing.T) {
	t.Run("type user.Metadata", func(t *testing.T) {
		reqMetadata := []*user.Metadata{
			{Key: "key1", Value: []byte("value1")},
			{Key: "key2", Value: []byte("value2")},
			{Key: "key3", Value: nil},
			{Key: "key4", Value: []byte("")},
		}
		want := []*domain.Metadata{
			{Key: "key1", Value: []byte("value1")},
			{Key: "key2", Value: []byte("value2")},
			{Key: "key3", Value: nil},
			{Key: "key4", Value: []byte("")},
		}
		got := setUserMetadataToDomain(reqMetadata)
		assert.Equal(t, want, got)
	})

	t.Run("type user.UpdateMetadata", func(t *testing.T) {
		reqMetadata := []*user.UpdateMetadata{
			{Key: "key5", Value: []byte("value5")},
			{Key: "key6", Value: []byte("value7")},
			{Key: "key7", Value: nil},
			{Key: "key8", Value: []byte("")},
		}
		want := []*domain.Metadata{
			{Key: "key5", Value: []byte("value5")},
			{Key: "key6", Value: []byte("value7")},
			{Key: "key7", Value: nil},
			{Key: "key8", Value: []byte("")},
		}
		got := setUserMetadataToDomain(reqMetadata)
		assert.Equal(t, want, got)
	})

	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, setUserMetadataToDomain([]*user.Metadata{}))
	})

	t.Run("nil", func(t *testing.T) {
		var reqMetadata []*user.Metadata
		assert.Nil(t, setUserMetadataToDomain(reqMetadata))
	})

}
