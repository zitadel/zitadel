package eventstore

import "testing"

func TestOwnerTag(t *testing.T) {
	tests := []struct {
		name string
		kind string
		id   string
		want string
	}{
		{name: "org", kind: UniqueConstraintOwnerOrg, id: "123", want: "org:123"},
		{name: "empty id", kind: UniqueConstraintOwnerOrg, id: "", want: ""},
		{name: "empty kind", kind: "", id: "123", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OwnerTag(tt.kind, tt.id); got != tt.want {
				t.Errorf("OwnerTag() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUniqueConstraintAction_Valid(t *testing.T) {
	if !UniqueConstraintAdd.Valid() {
		t.Error("UniqueConstraintAdd should be valid")
	}
	if !UniqueConstraintRemove.Valid() {
		t.Error("UniqueConstraintRemove should be valid")
	}
	if !UniqueConstraintInstanceRemove.Valid() {
		t.Error("UniqueConstraintInstanceRemove should be valid")
	}
	if !UniqueConstraintRemoveByOwner.Valid() {
		t.Error("UniqueConstraintRemoveByOwner should be valid")
	}
	if UniqueConstraintAction(uniqueConstraintActionCount).Valid() {
		t.Error("uniqueConstraintActionCount should be invalid")
	}
}

func TestUniqueConstraint_WithOwners(t *testing.T) {
	got := NewAddEventUniqueConstraint("usernames", "alice", "err").WithOwners(
		OwnerTag(UniqueConstraintOwnerOrg, "org1"),
		"",
		OwnerTag(UniqueConstraintOwnerUser, "user1"),
	)
	if len(got.Owners) != 2 || got.Owners[0] != "org:org1" || got.Owners[1] != "user:user1" {
		t.Errorf("unexpected owners: %#v", got.Owners)
	}
	if NewRemoveUniqueConstraintsByOwner(UniqueConstraintOwnerOrg, "org1").Action != UniqueConstraintRemoveByOwner {
		t.Error("expected UniqueConstraintRemoveByOwner")
	}
}
