package ldap

//go:generate mockgen -package mock -destination ./mock/ldap.mock.go github.com/zitadel/zitadel/internal/idp/providers/ldap ProviderInterface
//go:generate mockgen -package mock -destination ./mock/session.mock.go github.com/zitadel/zitadel/internal/idp Session
