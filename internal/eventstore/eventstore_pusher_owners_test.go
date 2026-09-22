package eventstore_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestEventstore_Push_UniqueConstraintOwners(t *testing.T) {
	for pusherName, pusher := range pushers {
		t.Run(pusherName, func(t *testing.T) {
			runUniqueConstraintOwnerCases(t, pusher, clients[pusherName])
		})
	}
}

func runUniqueConstraintOwnerCases(t *testing.T, pusher eventstore.Pusher, client *database.DB) {
	t.Helper()
	db := eventstore.NewEventstore(&eventstore.Config{
		Querier: queriers["v2(inmemory)"],
		Pusher:  pusher,
	})
	const instanceID = "owners-instance"
	for _, tc := range uniqueConstraintOwnerCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(cleanupEventstore(client))
			tc.run(t, db, client, instanceID)
		})
	}
}

type uniqueConstraintOwnerCase struct {
	name string
	run  func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string)
}

var uniqueConstraintOwnerCases = []uniqueConstraintOwnerCase{
	{
		name: "add stores owners",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			_, err := db.Push(context.Background(), generateCommand("owners-add", "1",
				withInstanceID(instanceID),
				generateAddUniqueConstraint("usernames", "alice", "org:org-1", "user:user-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertOwners(t, client, instanceID, "usernames", "alice", []string{"org:org-1", "user:user-1"})
		},
	},
	{
		name: "omit owners stores empty",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			_, err := db.Push(context.Background(), generateCommand("owners-empty", "1",
				withInstanceID(instanceID),
				generateAddUniqueConstraint("usernames", "bob"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertOwners(t, client, instanceID, "usernames", "bob", nil)
		},
	},
	{
		name: "add duplicate already exists",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			cmd := func() eventstore.Command {
				return generateCommand("owners-dup", "1",
					withInstanceID(instanceID),
					generateAddUniqueConstraint("usernames", "alice", "org:org-1"),
				)
			}
			if _, err := db.Push(context.Background(), cmd()); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), cmd())
			if !zerrors.IsErrorAlreadyExists(err) {
				t.Fatalf("expected already exists, got %v", err)
			}
		},
	},
	{
		name: "remove by field case insensitive",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "mixedcase", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-field-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraint("usernames", "MixedCase"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 0)
		},
	},
	{
		name: "remove by owner org",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "alice", "org:org-x", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "bob", "org:org-y", "user:u2"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-org-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "org-x"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "bob", []string{"org:org-y", "user:u2"})
		},
	},
	{
		name: "remove by owner user",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "alice", "org:org-1", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "external_idps", "idp1ext1", "org:org-1", "idp:idp1", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "bob", "org:org-1", "user:u2"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-user-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerUser, "u1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "bob", []string{"org:org-1", "user:u2"})
		},
	},
	{
		name: "remove by owner idp",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "external_idps", "idp1ext1", "org:org-1", "idp:idp-1", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "external_idps", "idp2ext2", "org:org-1", "idp:idp-2", "user:u2"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-idp-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerIDP, "idp-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "external_idps", "idp2ext2", []string{"org:org-1", "idp:idp-2", "user:u2"})
		},
	},
	{
		name: "two org tags either org deletes grant",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "project_grant", "granted:project", "org:granting", "org:granted", "project:project"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-grant-org", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "granted"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 0)
		},
	},
	{
		name: "instance remove deletes tagged rows",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "alice", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, "other", "usernames", "alice", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-instance-remove", "1",
				withInstanceID(instanceID),
				generateRemoveInstanceUniqueConstraints(),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 0)
			assertUniqueCount(t, client, "other", 1)
		},
	},
	{
		name: "global and empty owners not deleted by owner",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, "", "instance_domain", "example.com"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "untagged"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "tagged", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-skip-empty", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "org-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, "", 1)
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "untagged", nil)
		},
	},
	{
		name: "mixed add tagged and empty owners",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			_, err := db.Push(context.Background(), generateCommand("owners-mixed-add", "1",
				withInstanceID(instanceID),
				generateAddUniqueConstraint("usernames", "tagged", "org:org-1", "user:u1"),
				generateAddUniqueConstraint("usernames", "untagged"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 2)
			assertOwners(t, client, instanceID, "usernames", "tagged", []string{"org:org-1", "user:u1"})
			assertOwners(t, client, instanceID, "usernames", "untagged", nil)
		},
	},
	{
		name: "mixed batch deletes then inserts",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "old", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "org_name", "acme", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-mixed", "1",
				withInstanceID(instanceID),
				generateAddUniqueConstraint("usernames", "new", "org:org-1", "user:u1"),
				generateRemoveUniqueConstraint("usernames", "old"),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "org-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "new", []string{"org:org-1", "user:u1"})
		},
	},
}

func TestEventstore_Push_UniqueConstraintOwners_MissingColumn(t *testing.T) {
	for pusherName, pusher := range pushers {
		t.Run(pusherName, func(t *testing.T) {
			client := clients[pusherName]
			tx, err := client.Begin()
			if err != nil {
				t.Fatal(err)
			}
			// DROP COLUMN takes AccessExclusiveLock until rollback. Other tests
			// against the same database may stall; they still complete afterwards.
			t.Cleanup(func() {
				if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
					t.Errorf("rollback missing-column tx: %v", err)
				}
			})
			if _, err := tx.Exec("ALTER TABLE eventstore.unique_constraints DROP COLUMN IF EXISTS owners"); err != nil {
				t.Fatal(err)
			}

			db := eventstore.NewEventstore(&eventstore.Config{
				Querier: queriers["v2(inmemory)"],
				Pusher:  pusher,
			})
			const instanceID = "owners-missing-col"

			t.Run("empty owners add succeeds", func(t *testing.T) {
				_, err := db.PushWithClient(context.Background(), tx, generateCommand("owners-missing-empty", "1",
					withInstanceID(instanceID),
					generateAddUniqueConstraint("usernames", "bob"),
				))
				if err != nil {
					t.Fatal(err)
				}
				assertUniqueCountTx(t, tx, instanceID, 1)
			})

			t.Run("tagged owners add fails", func(t *testing.T) {
				_, err := db.PushWithClient(context.Background(), tx, generateCommand("owners-missing-tagged", "1",
					withInstanceID(instanceID),
					generateAddUniqueConstraint("usernames", "alice", "org:org-1", "user:u1"),
				))
				if err == nil {
					t.Fatal("expected missing owners column error")
				}
				if !isUndefinedColumn(err) {
					t.Fatalf("expected undefined column, got %v", err)
				}
			})
		})
	}
}

func isUndefinedColumn(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42703"
}

func assertUniqueCountTx(t *testing.T, tx *sql.Tx, instanceID string, want int) {
	t.Helper()
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM eventstore.unique_constraints WHERE instance_id = $1", instanceID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Errorf("unique count = %d, want %d", count, want)
	}
}

func assertUniqueCount(t *testing.T, db *database.DB, instanceID string, want int) {
	t.Helper()
	var count int
	err := db.QueryRow(func(row *sql.Row) error {
		return row.Scan(&count)
	}, "SELECT COUNT(*) FROM eventstore.unique_constraints WHERE instance_id = $1", instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Errorf("unique count = %d, want %d", count, want)
	}
}

func assertOwners(t *testing.T, db *database.DB, instanceID, uniqueType, uniqueField string, want []string) {
	t.Helper()
	var owners database.TextArray[string]
	err := db.QueryRow(func(row *sql.Row) error {
		return row.Scan(&owners)
	}, "SELECT owners FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3", instanceID, uniqueType, uniqueField)
	if err != nil {
		t.Fatal(err)
	}
	if len(owners) != len(want) {
		t.Fatalf("owners = %v, want %v", []string(owners), want)
	}
	got := map[string]int{}
	for _, owner := range owners {
		got[owner]++
	}
	for _, owner := range want {
		if got[owner] == 0 {
			t.Fatalf("owners = %v, want %v", []string(owners), want)
		}
		got[owner]--
	}
}
