// Package repository implements the repositories defined in the domain package.
// The repositories are used by the domain package to access the database.
// the inheritance.sql file is me over-engineering table inheritance.
// I would create a user table which is inherited by human_user and machine_user and the same for objects like idps.
package repository
