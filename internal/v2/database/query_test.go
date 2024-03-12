package database

// type Query struct {
// 	Columns  []Column
// 	From     From
// 	Join     Join
// 	Where    Where
// 	Limit    uint32
// 	Offset   uint32
// 	GroupBy  GroupBy
// 	Having   Having
// 	Ordering Ordering
// }

// type Column interface {
// 	Write(stmt *Statement)
// 	WriteIdentifier(stmt *Statement)
// }

// type column struct {
// 	Name string
// 	As   string
// }

// type function struct {
// 	Name   string
// 	Column Column
// }

// func (f *function) Write(stmt *Statement) {
// 	stmt.Builder.WriteString(f.Name)
// 	stmt.Builder.WriteRune('(')
// 	f.Column.WriteIdentifier(stmt)
// 	stmt.Builder.WriteRune(')')
// }

// type From interface {
// 	Write(stmt *Statement)
// }

// type Join interface {
// 	Write(stmt *Statement)
// }

// type Where interface {
// 	Write(stmt *Statement)
// }

// type where struct{}

// type condition struct{}

// type Ordering interface {
// 	Write(stmt *Statement)
// }

// type orderBy struct {
// 	Fields []*orderByField
// }

// func (ob *orderBy) Write(stmt *Statement) {
// 	for i, field := range ob.Fields {
// 		field.Write(stmt)
// 		if i < len(ob.Fields)-1 {
// 			stmt.Builder.WriteString(", ")
// 		}
// 	}
// }

// type orderByField struct {
// 	Identifier string
// 	Desc       bool
// }

// func (f *orderByField) Write(stmt *Statement) {
// 	stmt.Builder.WriteString(f.Identifier)
// 	if f.Desc {
// 		stmt.Builder.WriteString(" DESC")
// 	}
// }

// type GroupBy interface {
// 	Write(stmt *Statement)
// }

// type groupBy struct {
// 	Fields []*groupByField
// }

// func (gb *groupBy) Write(stmt *Statement) {
// 	for i, field := range gb.Fields {
// 		field.Write(stmt)
// 		if i < len(gb.Fields)-1 {
// 			stmt.Builder.WriteString(", ")
// 		}
// 	}
// }

// type groupByField struct {
// 	Identifier string
// }

// func (f *groupByField) Write(stmt *Statement) {
// 	stmt.Builder.WriteString(f.Identifier)
// }

// type Having interface {
// 	Write(stmt *Statement)
// }

// type from interface {
// 	~string | *Query
// }
