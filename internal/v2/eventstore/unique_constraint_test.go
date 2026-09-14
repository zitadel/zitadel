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
	if !UniqueConstraintRemoveByOwner.Valid() {
		t.Error("UniqueConstraintRemoveByOwner should be valid")
	}
	if UniqueConstraintAction(uniqueConstraintActionCount).Valid() {
		t.Error("uniqueConstraintActionCount should be invalid")
	}
}
