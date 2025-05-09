package stmt

// type table struct {
// 	schema string
// 	name   string

// 	possibleJoins []*join

// 	columns []*col
// }

// type col struct {
// 	*table

// 	name string
// }

// type join struct {
// 	*table

// 	on []*joinColumns
// }

// type joinColumns struct {
// 	left, right *col
// }

// var (
// 	userTable = &table{
// 		schema: "zitadel",
// 		name:   "users",
// 	}
// 	userColumns = []*col{
// 		userInstanceIDColumn,
// 		userOrgIDColumn,
// 		userIDColumn,
// 		userUsernameColumn,
// 	}
// 	userInstanceIDColumn = &col{
// 		table: userTable,
// 		name:  "instance_id",
// 	}
// 	userOrgIDColumn = &col{
// 		table: userTable,
// 		name:  "org_id",
// 	}
// 	userIDColumn = &col{
// 		table: userTable,
// 		name:  "id",
// 	}
// 	userUsernameColumn = &col{
// 		table: userTable,
// 		name:  "username",
// 	}
// 	userJoins = []*join{
// 		{
// 			table: instanceTable,
// 			on: []*joinColumns{
// 				{
// 					left:  instanceIDColumn,
// 					right: userInstanceIDColumn,
// 				},
// 			},
// 		},
// 		{
// 			table: orgTable,
// 			on: []*joinColumns{
// 				{
// 					left:  orgIDColumn,
// 					right: userOrgIDColumn,
// 				},
// 			},
// 		},
// 	}
// )

// var (
// 	instanceTable = &table{
// 		schema: "zitadel",
// 		name:   "instances",
// 	}
// 	instanceColumns = []*col{
// 		instanceIDColumn,
// 		instanceNameColumn,
// 	}
// 	instanceIDColumn = &col{
// 		table: instanceTable,
// 		name:  "id",
// 	}
// 	instanceNameColumn = &col{
// 		table: instanceTable,
// 		name:  "name",
// 	}
// )

// var (
// 	orgTable = &table{
// 		schema: "zitadel",
// 		name:   "orgs",
// 	}
// 	orgColumns = []*col{
// 		orgInstanceIDColumn,
// 		orgIDColumn,
// 		orgNameColumn,
// 	}
// 	orgInstanceIDColumn = &col{
// 		table: orgTable,
// 		name:  "instance_id",
// 	}
// 	orgIDColumn = &col{
// 		table: orgTable,
// 		name:  "id",
// 	}
// 	orgNameColumn = &col{
// 		table: orgTable,
// 		name:  "name",
// 	}
// )

// func init() {
// 	instanceTable.columns = instanceColumns
// 	userTable.columns = userColumns

// 	userTable.possibleJoins = []join{
// 		{
// 			table: userTable,
// 			on: []joinColumns{
// 				{
// 					left:  userIDColumn,
// 					right: userIDColumn,
// 				},
// 			},
// 		},
// 	}
// }
