package repository

//AggregateType is the object name
type AggregateType string

// //Aggregate represents an object
// type Aggregate struct {
// 	//ID id is the unique identifier of the aggregate
// 	// the client must generate it by it's own
// 	ID string
// 	//Type describes the meaning of this aggregate
// 	// it could an object like user
// 	Type AggregateType

// 	//ResourceOwner is the organisation which owns this aggregate
// 	// an aggregate can only be managed by one organisation
// 	// use the ID of the org
// 	ResourceOwner string

// 	//Events describe all the changes made on an aggregate
// 	Events []*Event
// }
