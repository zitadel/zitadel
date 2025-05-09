package storage

// import "context"

// type Storage interface {
// 	Write(ctx context.Context, commands ...Command) error
// }

// type Command interface {
// 	// Subjects returns a list of subjects which describe the command
// 	// the type should always be in past tense
// 	// "." is used as separator
// 	// e.g. "user.created"
// 	Type() string
// 	// Payload returns the payload of the command
// 	// The payload can either be
// 	// - a struct
// 	// - a map
// 	// - nil
// 	Payload() any
// 	// Revision returns the revision of the command
// 	Revision() uint8
// 	// Creator returns the user id who created the command
// 	Creator() string

// 	// Object returns the object the command belongs to
// 	Object() Model
// 	// Parents returns the parents of the object
// 	// If the list is empty there are no parents
// 	Parents() []Model

// 	// Objects returns the models to update during inserting the commands
// 	Objects() []Object
// }

// type Model struct {
// 	Name string
// 	ID   string
// }

// type Object struct {
// 	Model Model
// }
