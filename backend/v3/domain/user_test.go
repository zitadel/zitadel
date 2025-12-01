package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
)

func TestPasskeyList_GetPasskeysOfType(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		passkeys    domain.PasskeyList
		wantedTypes []domain.PasskeyType
		want        domain.PasskeyList
	}{
		{
			name:        "empty passkey list",
			passkeys:    domain.PasskeyList{},
			wantedTypes: []domain.PasskeyType{domain.PasskeyTypePasswordless},
			want:        domain.PasskeyList{},
		},
		{
			name: "single matching type",
			passkeys: domain.PasskeyList{
				{ID: "1", Type: domain.PasskeyTypePasswordless},
				{ID: "2", Type: domain.PasskeyTypeU2F},
			},
			wantedTypes: []domain.PasskeyType{domain.PasskeyTypePasswordless},
			want: domain.PasskeyList{
				{ID: "1", Type: domain.PasskeyTypePasswordless},
			},
		},
		{
			name: "multiple matching types",
			passkeys: domain.PasskeyList{
				{ID: "1", Type: domain.PasskeyTypePasswordless},
				{ID: "2", Type: domain.PasskeyTypeU2F},
				{ID: "3", Type: domain.PasskeyTypePasswordless},
			},
			wantedTypes: []domain.PasskeyType{domain.PasskeyTypePasswordless, domain.PasskeyTypeU2F},
			want: domain.PasskeyList{
				{ID: "1", Type: domain.PasskeyTypePasswordless},
				{ID: "2", Type: domain.PasskeyTypeU2F},
				{ID: "3", Type: domain.PasskeyTypePasswordless},
			},
		},
		{
			name: "no matching types",
			passkeys: domain.PasskeyList{
				{ID: "1", Type: domain.PasskeyTypePasswordless},
				{ID: "2", Type: domain.PasskeyTypeU2F},
			},
			wantedTypes: []domain.PasskeyType{domain.PasskeyTypeUnspecified},
			want:        domain.PasskeyList{},
		},
		{
			name: "empty wanted types",
			passkeys: domain.PasskeyList{
				{ID: "1", Type: domain.PasskeyTypePasswordless},
				{ID: "2", Type: domain.PasskeyTypeU2F},
			},
			wantedTypes: []domain.PasskeyType{},
			want:        domain.PasskeyList{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.passkeys.GetPasskeysOfType(tc.wantedTypes)
			require.Len(t, got, len(tc.want))
			for i, p := range got {
				assert.Equal(t, tc.want[i].ID, p.ID)
				assert.Equal(t, tc.want[i].Type, p.Type)
			}
		})
	}
}
