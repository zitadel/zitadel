package service

type DomainGenerator interface {
	GenerateDomain() (string, error)
}

type Domain struct {
	Domain      string
	IsPrimary   bool
	IsGenerated bool
}
