package domain

type UniqueConstraintMigration struct {
	AggregateID  string
	ObjectID     string
	UniqueType   string
	UniqueField  string
	ErrorMessage string
}
