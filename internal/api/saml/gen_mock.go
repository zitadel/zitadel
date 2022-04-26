package saml

//go:generate mockgen -package mock -destination ./mock/storage.mock.go github.com/caos/zitadel/internal/api/saml Storage
//go:generate mockgen -package mock -destination ./mock/idpstorage.mock.go github.com/caos/zitadel/internal/api/saml IDPStorage
//go:generate mockgen -package mock -destination ./mock/authrequestint.mock.go github.com/caos/zitadel/internal/api/saml/models AuthRequestInt
