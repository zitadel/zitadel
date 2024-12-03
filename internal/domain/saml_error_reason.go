package domain

type SAMLErrorReason int32

const (
	SAMLErrorReasonUnspecified SAMLErrorReason = iota
)

func SAMLErrorReasonFromError(err error) SAMLErrorReason {
	return SAMLErrorReasonUnspecified
}
