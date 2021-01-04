//go generate
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

func main() {
	errorName := readErrorName()
	errorName = validateErrorName(errorName)

	data := &Data{
		ErrorName: errorName,
	}

	errorFile := data.createFile("error.go.tmpl")
	data.createTemplate("error.go.tmpl", errorFile)
	if err := errorFile.Close(); err != nil {
		log.Fatal(err)
	}

	testFile := data.createFile("error_test.go.tmpl")
	data.createTemplate("error_test.go.tmpl", testFile)
	if err := testFile.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Print(`
!!!!!
  Add status mapping in internal/api/grpc/caos_errors.go	
!!!!!`)
}

type Data struct {
	ErrorName string
}

func (data *Data) createFile(tmplName string) *os.File {
	filename := strings.Replace(tmplName, "error", strings.ToLower(data.ErrorName), 1)
	filename = filename[:len(filename)-5]
	filePath := fmt.Sprintf("../%s", filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("unable to create file (%s): %v", filePath, err)
	}
	return file
}

func (data *Data) createTemplate(templateName string, file *os.File) {
	tmpl := template.Must(template.New(templateName).ParseFiles(templateName))
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatal("unable to execute tmpl: ", err)
	}
}

func readErrorName() (errorName string) {
	flag.StringVar(&errorName, "Name", "", "KeyType of the error (e.g. Internal)")
	flag.Parse()
	return errorName
}

func validateErrorName(errorName string) string {
	if errorName == "" {
		log.Fatal("pass argument name")
	}
	if strings.Contains(errorName, " ") || strings.Contains(errorName, ".") {
		log.Fatal("name cannot contain spaces or points")
	}
	return strings.Title(errorName)
}
