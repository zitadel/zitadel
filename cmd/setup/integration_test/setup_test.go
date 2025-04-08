// Package setup_test implements tests for procedural PostgreSQL functions,
// created in the database during Zitadel setup.
// Tests depend on `zitadel setup` being run first and therefore is run as integration tests.
// A PGX connection is used directly to the integration test database.
// This package assumes the database server available as per integration test defaults.
// See the [ConnString] constant.
package setup_test
