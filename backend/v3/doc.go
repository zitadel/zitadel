// the test used the manly relies on the following patterns:
// - api:
//   - some example stubs for the grpc api, it maps the calls and responses to the domain objects
//
// - domain:
//   - hexagonal architecture, it defines its dependencies as interfaces and the dependencies must use the objects defined by this package
//   - command pattern which implements the changes
//   - the invoker decorates the commands by checking for events, tracing, logging, potentially caching, etc.
//   - the database connections are manged in this package
//   - the database connections are passed to the repositories
//
// - storage:
//   - repository pattern, the repositories are defined as interfaces and the implementations are in the storage package
//   - the repositories are used by the domain package to access the database
//   - the eventstore to store events. At the beginning it writes to the same events table as the /internal package, afterwards it writes to a different table
//
// - telemetry:
//   - logging for standard output
//   - tracing for distributed tracing
//   - metrics for monitoring
package v3
