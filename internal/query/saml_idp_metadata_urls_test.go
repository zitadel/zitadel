package query

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestQueries_ListSAMLIDPMetadataURLs(t *testing.T) {
	columns := []string{"idp_id", "instance_id", "resource_owner", "owner_type", "metadata_url"}
	tests := []struct {
		name       string
		expects    func(sqlmock.Sqlmock)
		wantResult []SAMLIDPMetadataURL
		wantErr    bool
	}{
		{
			name: "query error",
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(samlIDPMetadataURLsListQuery)).
					WithArgs(domain.IDPTypeSAML, domain.IDPStateActive, "instance-1", "idp-1", uint32(10)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(samlIDPMetadataURLsListQuery)).
					WithArgs(domain.IDPTypeSAML, domain.IDPStateActive, "instance-1", "idp-1", uint32(10)).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow("idp-2", "instance-1", "org-1", domain.IdentityProviderTypeOrg, "https://idp.example.com/metadata").
						AddRow("idp-1", "instance-2", "instance-2", domain.IdentityProviderTypeSystem, "https://other.example.com/metadata"))
			},
			wantResult: []SAMLIDPMetadataURL{
				{ID: "idp-2", InstanceID: "instance-1", ResourceOwner: "org-1", OwnerType: domain.IdentityProviderTypeOrg, MetadataURL: "https://idp.example.com/metadata"},
				{ID: "idp-1", InstanceID: "instance-2", ResourceOwner: "instance-2", OwnerType: domain.IdentityProviderTypeSystem, MetadataURL: "https://other.example.com/metadata"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			tt.expects(mock)

			q := &Queries{client: &database.DB{DB: db}}
			got, err := q.ListSAMLIDPMetadataURLs(context.Background(), "instance-1", "idp-1", 10)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantResult, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
