package rsync

import (
	"bytes"
	"text/template"
)

func GetConfigPartS3(
	name string,
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
) (string, error) {
	type file struct {
		Name            string
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
	}
	fileStruct := file{
		Name:            name,
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	tmpl, err := template.New(name).Parse(
		"[{{.Name}}]\n" +
			"type = s3\n" +
			"provider = Minio\n" +
			"env_auth = false\n" +
			"access_key_id = {{.AccessKeyID}}\n" +
			"secret_access_key = {{.SecretAccessKey}}\n" +
			"endpoint = {{.Endpoint}}\n" +
			"acl = private\n")
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
