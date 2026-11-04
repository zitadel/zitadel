package command

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/domain"
)

// LabelPolicyLinkPreviewUpdate holds the link preview text of a label policy change.
// A nil field keeps the current value, so clients that don't know the fields don't clear them.
type LabelPolicyLinkPreviewUpdate struct {
	Title       *string
	Description *string
}

func (u LabelPolicyLinkPreviewUpdate) validate() error {
	return domain.ValidateLinkPreviewText(gu.Value(u.Title), gu.Value(u.Description))
}
