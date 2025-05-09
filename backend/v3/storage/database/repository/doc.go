// Repository package provides the database repository for the application.
// It contains the implementation of the [repository pattern](https://martinfowler.com/eaaCatalog/repository.html) for the database.

// funcs which need to interact with the database should create interfaces which are implemented by the
// [query] and [exec] structs respectively their factory methods [Query] and [Execute]. The [query] struct is used for read operations, while the [exec] struct is used for write operations.

package repository
