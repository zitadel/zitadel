package query

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query/projection"
)

type Key interface {
	ID() string
	Algorithm() string
	Use() crypto.KeyUsage
	Sequence() uint64
}

type PrivateKey interface {
	Key
	Expiry() time.Time
	Key() *crypto.CryptoValue
}

type PublicKey interface {
	Key
	Expiry() time.Time
	Key() interface{}
}

type PublicKeys struct {
	SearchResponse
	Keys []PublicKey
}

type key struct {
	id            string
	creationDate  time.Time
	changeDate    time.Time
	sequence      uint64
	resourceOwner string
	algorithm     string
	use           crypto.KeyUsage
}

func (k *key) ID() string {
	return k.id
}

func (k *key) Algorithm() string {
	return k.algorithm
}

func (k *key) Use() crypto.KeyUsage {
	return k.use
}

func (k *key) Sequence() uint64 {
	return k.sequence
}

var (
	keyTable = table{
		name:          projection.KeyProjectionTable,
		instanceIDCol: projection.KeyColumnInstanceID,
	}
	KeyColID = Column{
		name:  projection.KeyColumnID,
		table: keyTable,
	}
	KeyColCreationDate = Column{
		name:  projection.KeyColumnCreationDate,
		table: keyTable,
	}
	KeyColChangeDate = Column{
		name:  projection.KeyColumnChangeDate,
		table: keyTable,
	}
	KeyColResourceOwner = Column{
		name:  projection.KeyColumnResourceOwner,
		table: keyTable,
	}
	KeyColInstanceID = Column{
		name:  projection.KeyColumnInstanceID,
		table: keyTable,
	}
	KeyColSequence = Column{
		name:  projection.KeyColumnSequence,
		table: keyTable,
	}
	KeyColAlgorithm = Column{
		name:  projection.KeyColumnAlgorithm,
		table: keyTable,
	}
	KeyColUse = Column{
		name:  projection.KeyColumnUse,
		table: keyTable,
	}
)

var (
	keyPrivateTable = table{
		name:          projection.KeyPrivateTable,
		instanceIDCol: projection.KeyPrivateColumnInstanceID,
	}
	KeyPrivateColID = Column{
		name:  projection.KeyPrivateColumnID,
		table: keyPrivateTable,
	}
	KeyPrivateColExpiry = Column{
		name:  projection.KeyPrivateColumnExpiry,
		table: keyPrivateTable,
	}
	KeyPrivateColKey = Column{
		name:  projection.KeyPrivateColumnKey,
		table: keyPrivateTable,
	}
)
