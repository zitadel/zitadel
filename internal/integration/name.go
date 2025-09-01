package integration

import "github.com/brianvoe/gofakeit/v6"

// company private function to add a random string to the gofakeit.Company function
func company() string {
	return gofakeit.Company() + "-" + RandString(5)
}

func OrganizationName() string {
	return company()
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
