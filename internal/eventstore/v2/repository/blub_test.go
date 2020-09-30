package repository

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
)

func TestBlub(t *testing.T) {
	s, db, kvDB := serverutils.StartServer(t, base.TestServerArgs{})
	defer s.Stopper().Stop()
	// If really needed, in tests that can depend on server, downcast to
	// server.TestServer:
	ts := s.(*server.TestServer)
}
