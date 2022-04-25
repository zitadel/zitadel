package saml

//go:generate mockgen -package mock -destination ./mock/generator.mock.go github.com/caos/zitadel/internal/api/saml Storage
//go:generate mockgen -package mock -destination ./mock/generator.mock.go github.com/caos/zitadel/internal/api/saml IDPStorage
