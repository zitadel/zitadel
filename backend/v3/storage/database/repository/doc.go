// Package repository implements the repositories defined in the domain package.
// The repositories are used by the domain package to access the database.
//
// Repository methods may require a minimal set of columns in the [database.Condition] to be passed.
// This is to ensure that the repository methods are used correctly and do not affect more rows than intended.
// For example:
//   - The instance ID is almost always required to ensure that only rows of the correct instance are affected.
//   - Update and Delete methods often require the columns from the primary key to ensure that only one row is affected.
// The required columns are checked during the call and an error is returned if they are not present.

// the inheritance.sql file is me(silvan) over-engineering table inheritance.
// I would create a user table which is inherited by human_user and machine_user and the same for objects like idps.
package repository
