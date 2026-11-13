package instance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
)

func TestNewSMTPConfigPasswordChangedEvent(t *testing.T) {
	password := &crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "key-id", Crypted: []byte("password")}
	event := NewSMTPConfigPasswordChangedEvent(context.Background(), &NewAggregate("instance-id").Aggregate, "config-id", password)

	assert.Equal(t, "config-id", event.ID)
	assert.Equal(t, password, event.Password)
}
