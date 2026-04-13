package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func TestDomainLinksModelToGRPCResponse(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		tcs := []struct {
			name     string
			input    []domain.Link
			expected []*settings.Link
		}{
			{
				name:     "empty list",
				input:    make([]domain.Link, 0),
				expected: make([]*settings.Link, 0),
			},
			{
				name: "normal case",
				input: []domain.Link{
					{Type: domain.LinkTypeDocs, Url: "https://docs.example.com"},
					{Type: domain.LinkTypePrivacyPolicy, Url: "https://privacy.example.com", Target: domain.LinkTargetBlank},
					{Type: domain.LinkTypeTermsOfService, Url: "https://terms.example.com", Target: domain.LinkTargetSelf},
					{Type: domain.LinkTypeHelp, Url: "https://help.example.com", Target: domain.LinkTargetBlank},
					{Type: domain.LinkTypeSupport, Url: "mailto:support@example.com"},
					{
						Type:           domain.LinkTypeCustom,
						Url:            "https://shop.example.com",
						Target:         domain.LinkTargetSelf,
						TranslationKey: "shop",
					},
				},
				expected: []*settings.Link{
					{Type: settings.LinkType_LINK_TYPE_DOCS, Url: "https://docs.example.com", Target: settings.LinkTarget_LINK_TARGET_UNSPECIFIED},
					{Type: settings.LinkType_LINK_TYPE_PRIVACY_POLICY, Url: "https://privacy.example.com", Target: settings.LinkTarget_LINK_TARGET_BLANK},
					{Type: settings.LinkType_LINK_TYPE_TERMS_OF_SERVICE, Url: "https://terms.example.com", Target: settings.LinkTarget_LINK_TARGET_SELF},
					{Type: settings.LinkType_LINK_TYPE_HELP, Url: "https://help.example.com", Target: settings.LinkTarget_LINK_TARGET_BLANK},
					{Type: settings.LinkType_LINK_TYPE_SUPPORT, Url: "mailto:support@example.com"},
					{
						Type:           settings.LinkType_LINK_TYPE_CUSTOM,
						Url:            "https://shop.example.com",
						Target:         settings.LinkTarget_LINK_TARGET_SELF,
						TranslationKey: "shop",
					},
				},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				val, err := DomainLinksModelToGRPCResponse(tc.input)
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expected, val)
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		tcs := []struct {
			name        string
			input       []domain.Link
			expectedErr error
		}{
			{
				name:        "unknown link type",
				input:       []domain.Link{{Type: 999}},
				expectedErr: zerrors.ThrowInvalidArgumentf(nil, "BCgEF3", "unknown link type %v", 999),
			},
			{
				name:        "unknown link target",
				input:       []domain.Link{{Target: 999}},
				expectedErr: zerrors.ThrowInvalidArgumentf(nil, "W35qZF", "unknown target type %v", 999),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				val, err := DomainLinksModelToGRPCResponse(tc.input)
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.ErrorIs(t, tc.expectedErr, err)
			})
		}
	})
}

func TestGrpcLinksToDomain(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		tcs := []struct {
			name     string
			input    []*settings.Link
			expected []domain.Link
		}{
			{
				name:     "empty list",
				input:    make([]*settings.Link, 0),
				expected: make([]domain.Link, 0),
			},
			{
				name: "normal case",
				input: []*settings.Link{
					{Type: settings.LinkType_LINK_TYPE_DOCS, Url: "https://docs.example.com", Target: settings.LinkTarget_LINK_TARGET_UNSPECIFIED},
					{Type: settings.LinkType_LINK_TYPE_PRIVACY_POLICY, Url: "https://privacy.example.com", Target: settings.LinkTarget_LINK_TARGET_BLANK},
					{Type: settings.LinkType_LINK_TYPE_TERMS_OF_SERVICE, Url: "https://terms.example.com", Target: settings.LinkTarget_LINK_TARGET_SELF},
					{Type: settings.LinkType_LINK_TYPE_HELP, Url: "https://help.example.com", Target: settings.LinkTarget_LINK_TARGET_BLANK},
					{Type: settings.LinkType_LINK_TYPE_SUPPORT, Url: "mailto:support@example.com"},
					{
						Type:           settings.LinkType_LINK_TYPE_CUSTOM,
						Url:            "https://shop.example.com",
						Target:         settings.LinkTarget_LINK_TARGET_SELF,
						TranslationKey: "shop",
					},
				},
				expected: []domain.Link{
					{Type: domain.LinkTypeDocs, Url: "https://docs.example.com"},
					{Type: domain.LinkTypePrivacyPolicy, Url: "https://privacy.example.com", Target: domain.LinkTargetBlank},
					{Type: domain.LinkTypeTermsOfService, Url: "https://terms.example.com", Target: domain.LinkTargetSelf},
					{Type: domain.LinkTypeHelp, Url: "https://help.example.com", Target: domain.LinkTargetBlank},
					{Type: domain.LinkTypeSupport, Url: "mailto:support@example.com"},
					{
						Type:           domain.LinkTypeCustom,
						Url:            "https://shop.example.com",
						Target:         domain.LinkTargetSelf,
						TranslationKey: "shop",
					},
				},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				val, err := GrpcLinksToDomain(tc.input)
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expected, val)
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		tcs := []struct {
			name        string
			input       []*settings.Link
			expectedErr error
		}{
			{
				name:        "unknown link type",
				input:       []*settings.Link{{Type: 999}},
				expectedErr: zerrors.ThrowInvalidArgumentf(nil, "GJeupd", "unknown link type %v", 999),
			},
			{
				name:        "unknown link target",
				input:       []*settings.Link{{Target: 999}},
				expectedErr: zerrors.ThrowInvalidArgumentf(nil, "yJJTA3", "unknown target type %v", 999),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				val, err := GrpcLinksToDomain(tc.input)
				assert.Nil(t, val)
				assert.Error(t, err)
				assert.ErrorIs(t, tc.expectedErr, err)
			})
		}
	})
}
