// This package contains the migration logic for the PostgreSQL dialect.
// It uses the [github.com/jackc/tern/v2/migrate] package to handle the migration process.
//
// **Developer Note**:
//
// Each migration MUST be registered in an init function.
// Create a go file for each migration with the sequence of the migration as prefix and some descriptive name.
// The file name MUST be in the format <sequence>_<name>.go.
// Each migration SHOULD provide an up and down migration.
// Prefer to write SQL statements instead of funcs if it is reasonable.
package migration
