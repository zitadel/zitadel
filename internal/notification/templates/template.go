package templates

import (
	"bytes"
	"html/template"
)

const (
	templateFileName = "template.html"
)

func ParseTemplate(fileName string, data interface{}) (string, error) {
	if fileName == "" {
		fileName = templateFileName
	}
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
