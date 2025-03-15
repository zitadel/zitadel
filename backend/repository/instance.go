package repository

type Instance struct {
	ID   string
	Name string
}

type ListRequest struct {
	Limit uint16
}
