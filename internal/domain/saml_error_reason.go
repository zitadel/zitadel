package domain

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

func SAMLErrorReasonFromError(err error) SAMLErrorReason {
	return SAMLErrorReasonUnspecified
}
