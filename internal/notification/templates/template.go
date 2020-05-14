package templates

import (
	"bytes"
	"html/template"
)

const (
	templateFileName = "template.html"
)

func GetParsedTemplate(contentData interface{}) (string, error) {
	template, err := ParseTemplateFile("", contentData)
	if err != nil {
		return "", err
	}
	return ParseTemplateText(template, contentData)
}

func ParseTemplateFile(fileName string, data interface{}) (string, error) {
	if fileName == "" {
		fileName = templateFileName
	}
	template, err := template.ParseFiles(fileName)
	if err != nil {
		return "", err
	}
	return parseTemplate(template, data)
}

func ParseTemplateText(text string, data interface{}) (string, error) {
	template, err := template.New("template").Parse(text)
	if err != nil {
		return "", err
	}
	return parseTemplate(template, data)
}

func parseTemplate(template *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := template.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
