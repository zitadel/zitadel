## JWT Profile – key loading
internal/api/oidc/integration_test/token_jwt_profile_test.go:118:			tokenSource, err := profile.NewJWTProfileTokenSourceFromKeyFileData(CTX, Instance.OIDCIssuer(), tt.keyData, tt.scope)
internal/api/oidc/token_exchange.go:141:		verifier := op.NewJWTProfileVerifierKeySet(keySetMap(client.client.PublicKeys), op.IssuerFromContext(ctx), time.Hour, client.client.ClockSkew, s.jwtProfileUserCheck(ctx, &resourceOwner, &preferredLanguage))
internal/api/oidc/introspect.go:157:			verifier := op.NewJWTProfileVerifierKeySet(keySetMap(client.PublicKeys), op.IssuerFromContext(ctx), time.Hour, time.Second)
internal/api/oidc/client.go:1022:	verifier := op.NewJWTProfileVerifierKeySet(keySetMap(client.PublicKeys), op.IssuerFromContext(ctx), time.Hour, client.ClockSkew)
internal/api/oidc/token_jwt_profile.go:73:	verifier := op.NewJWTProfileVerifier(
internal/api/authz/system_token.go:44:		systemJWTProfile: op.NewJWTProfileVerifier(
internal/api/oidc/token_exchange.go:141:		verifier := op.NewJWTProfileVerifierKeySet(keySetMap(client.client.PublicKeys), op.IssuerFromContext(ctx), time.Hour, client.client.ClockSkew, s.jwtProfileUserCheck(ctx, &resourceOwner, &preferredLanguage))

## startCaches
internal/api/ui/login/login.go:221:func startCaches(background context.Context, connectors connector.Connectors, federateLogoutCache cache.Cache[federatedlogout.Index, string, *federatedlogout.FederatedLogout]) (_ *Caches, err error) {
internal/command/cache.go:14:func startCaches(background context.Context, connectors connector.Connectors) (_ *Caches, err error) {
internal/query/cache.go:27:func startCaches(background context.Context, connectors connector.Connectors, instanceConfig ActiveInstanceConfig) (_ *Caches, err error) {

## RegisterCacheInvalidation
internal/eventstore/handler/v2/handler.go:444:// RegisterCacheInvalidation registers a function to be called when a cache needs to be invalidated.
internal/eventstore/handler/v2/handler.go:446:func (h *Handler) RegisterCacheInvalidation(invalidate func(ctx context.Context, aggregates []*eventstore.Aggregate)) {
internal/query/org.go:546:	projection.OrgProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:609:	projection.InstanceProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:610:	projection.InstanceDomainProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:611:	projection.InstanceFeatureProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:612:	projection.InstanceTrustedDomainProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:613:	projection.SecurityPolicyProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:617:	projection.LimitsProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:618:	projection.RestrictionsProjection.RegisterCacheInvalidation(invalidate)
internal/query/instance.go:621:	projection.SystemFeatureProjection.RegisterCacheInvalidation(func(ctx context.Context, _ []*eventstore.Aggregate) {
