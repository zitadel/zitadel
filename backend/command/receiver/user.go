package receiver

type User struct {
	ID       string
	Username string

	Email *Email
	Phone *Phone
}
