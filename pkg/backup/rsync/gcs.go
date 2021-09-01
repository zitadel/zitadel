package rsync

import (
	"bytes"
	"text/template"
)

func GetConfigPartGCS(
	name string,
	saSecretPath string,
) (string, error) {
	type file struct {
		Name                   string
		Endpoint               string
		ServiceAccountFilePath string
	}
	fileStruct := file{
		Name:                   name,
		ServiceAccountFilePath: saSecretPath,
	}
	tmpl, err := template.New(name).Parse(
		"[{{.Name}}]\n" +
			"type = google cloud storage\n" +
			"service_account_file = {{.ServiceAccountFilePath}}\n" +
			"object_acl = private\n" +
			"bucket_acl = private\n" +
			"location = europe-west")
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, fileStruct)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
