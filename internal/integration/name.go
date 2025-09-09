package integration

import "github.com/brianvoe/gofakeit/v6"

// company private function to add a random string to the gofakeit.Company function
func company() string {
	return gofakeit.Company() + "-" + RandString(5)
}

func OrganizationName() string {
	return company()
}

func Email() string {
	return RandString(5) + gofakeit.Email()
}

func Phone() string {
	return "+41" + gofakeit.Phone()
}

func FirstName() string {
	return gofakeit.FirstName()
}

func LastName() string {
	return gofakeit.LastName()
}

func Username() string {
	return gofakeit.Username() + RandString(5)
}

func Language() string {
	return gofakeit.LanguageBCP()
}

func UserSchemaName() string {
	return gofakeit.Name() + RandString(5)
}

// appName private function to add a random string to the gofakeit.AppName function
func appName() string {
	return gofakeit.AppName() + "-" + RandString(5)
}

func TargetName() string {
	return appName()
}

func ApplicationName() string {
	return appName()
}

func ProjectName() string {
	return appName()
}

func IDPName() string {
	return appName()
}

func RoleKey() string {
	return appName()
}

func RoleDisplayName() string {
	return appName()
}

func DomainName() string {
	return RandString(5) + gofakeit.DomainName()
}

func URL() string {
	return gofakeit.URL()
}

func RelayState() string {
	return ID()
}

func ID() string {
	return RandString(20)
}

func Slogan() string {
	return gofakeit.Slogan()
}

func Number() int {
	return gofakeit.Number(0, 1_000_000)
}
