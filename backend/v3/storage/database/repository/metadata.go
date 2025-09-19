package repository

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/backend/v3/storage/database"
// )

// // -------------------------------------------------------------
// // repository
// // -------------------------------------------------------------

// type metadata[
// 	PARENT any,
// 	RESULT any,
// ] struct {
// 	table string

// 	repository
// 	parent PARENT
// }

// // Get implements [domain.OrganizationMetadataRepository].
// func (o *metadata[PARENT, RESULT]) Get(ctx context.Context, opts ...database.QueryOption) (*RESULT, error) {
// 	panic("unimplemented")
// }

// // List implements [domain.OrganizationMetadataRepository].
// func (o *metadata[PARENT, RESULT]) List(ctx context.Context, opts ...database.QueryOption) ([]*RESULT, error) {
// 	panic("unimplemented")
// }

// // Set implements [domain.OrganizationMetadataRepository].
// func (o *metadata[PARENT, RESULT]) Set(ctx context.Context, metadata ...*RESULT) (int64, error) {
// 	panic("unimplemented")
// }

// // Remove implements [domain.OrganizationMetadataRepository].
// func (o *metadata[PARENT, RESULT]) Remove(ctx context.Context, condition database.Condition) (int64, error) {
// 	panic("unimplemented")
// }

// // -------------------------------------------------------------
// // conditions
// // -------------------------------------------------------------

// // InstanceIDCondition implements [domain.OrganizationMetadataRepository].
// func (o *orgMetadata) InstanceIDCondition(instanceID string) database.Condition {
// 	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
// }

// // KeyCondition implements [domain.OrganizationMetadataRepository].
// func (o *orgMetadata) KeyCondition(op database.TextOperation, key string) database.Condition {
// 	return database.NewTextCondition(o.KeyColumn(), op, key)
// }

// // ValueCondition implements [domain.OrganizationMetadataRepository].
// func (o *orgMetadata) ValueCondition(value any) database.Condition {
// 	// return database.NewValueCondition(o.ValueColumn(), value)
// 	panic("unimplemented")
// }
