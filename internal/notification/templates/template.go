package templates

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
)

const (
	templatesPath    = "/templates"
	templateFileName = "template.html"
)

func GetParsedTemplate(dir http.FileSystem, contentData interface{}) (string, error) {
	template, err := ParseTemplateFile(dir, "", contentData)
	if err != nil {
		return "", err
	}
	return ParseTemplateText(template, contentData)
}

func ParseTemplateFile(dir http.FileSystem, fileName string, data interface{}) (string, error) {
	if fileName == "" {
		fileName = templateFileName
	}
	template, err := readFile(dir, fileName)
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

func readFile(dir http.FileSystem, fileName string) (*template.Template, error) {
	f, err := dir.Open(templatesPath + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(fileName).Parse(string(content))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
