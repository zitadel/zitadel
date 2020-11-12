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

func GetParsedTemplate(mailhtml string, contentData interface{}) (string, error) {
	template, err := ParseTemplateFile(mailhtml, contentData)
	if err != nil {
		return "", err
	}
	return ParseTemplateText(template, contentData)
}

func ParseTemplateFile(mailhtml string, data interface{}) (string, error) {
	tmpl, err := template.New("tmpl").Parse(mailhtml)
	if err != nil {
		return "", err
	}

	return parseTemplate(tmpl, data)
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

func readFileFromDatabase(dir http.FileSystem, fileName string) (*template.Template, error) {
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
