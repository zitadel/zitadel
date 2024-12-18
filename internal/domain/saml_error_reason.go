package domain

import (
	"github.com/zitadel/saml/pkg/provider"
)

type SAMLErrorReason int32

const (
	SAMLErrorReasonUnspecified SAMLErrorReason = iota
	SAMLErrorReasonVersionMissmatch
	SAMLErrorReasonAuthNFailed
	SAMLErrorReasonInvalidAttrNameOrValue
	SAMLErrorReasonInvalidNameIDPolicy
	SAMLErrorReasonRequestDenied
	SAMLErrorReasonRequestUnsupported
	SAMLErrorReasonUnsupportedBinding
)

func SAMLErrorReasonToString(reason SAMLErrorReason) string {
	switch reason {
	case SAMLErrorReasonUnspecified:
		return "unspecified error"
	case SAMLErrorReasonVersionMissmatch:
		return provider.StatusCodeVersionMissmatch
	case SAMLErrorReasonAuthNFailed:
		return provider.StatusCodeAuthNFailed
	case SAMLErrorReasonInvalidAttrNameOrValue:
		return provider.StatusCodeInvalidAttrNameOrValue
	case SAMLErrorReasonInvalidNameIDPolicy:
		return provider.StatusCodeInvalidNameIDPolicy
	case SAMLErrorReasonRequestDenied:
		return provider.StatusCodeRequestDenied
	case SAMLErrorReasonRequestUnsupported:
		return provider.StatusCodeRequestUnsupported
	case SAMLErrorReasonUnsupportedBinding:
		return provider.StatusCodeUnsupportedBinding
	default:
		return "unspecified error"
	}
}
