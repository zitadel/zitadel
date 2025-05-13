// the test used the manly relies on the following patterns:
// - domain:
//   - hexagonal architecture, it defines its dependencies as interfaces and the dependencies must use the objects defined by this package
//   - command pattern which implements the changes
//   - the invoker decorates the commands by checking for events and tracing
//   - the database connections are manged in this package
//   - the database connections are passed to the repositories
//
// - storage:
//   - repository pattern, the repositories are defined as interfaces and the implementations are in the storage package
//   - the repositories are used by the domain package to access the database
package v3
