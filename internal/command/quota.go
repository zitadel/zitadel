package command

type QuotaSetup struct {
	Name         string
	CustomDomain string
	Human        AddHuman
	Roles        []string
}
