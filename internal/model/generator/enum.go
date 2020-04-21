//go generate
package generator

// import (
// 	"log"
// )

// /*
// type SearchMethod Enum

// var searchMethods = []string{"Equals", "StartsWith", "Contains"}
// */
// const tmpl = `

// type searchMethod int32

// func (s searchMethod) String() string {
// 	return searchMethods[s]
// }

// const (
// 	Equals searchMethod = 1 + iota
// 	StartsWith
// 	Contains
// )

// func SearchMethodToInt(s SearchMethod) int32 {
// 	return int32(s.(searchMethod))
// }

// func SearchMethodFromInt(index int32) SearchMethod {
// 	return searchMethod(index)
// }
// `

// func main() {
// 	errorName := readErrorName()
// 	errorName = validateErrorName(errorName)

// 	data := &Data{
// 		ErrorName: errorName,
// 	}

// 	errorFile := data.createFile("error.go.tmpl")
// 	data.createTemplate("error.go.tmpl", errorFile)
// 	if err := errorFile.Close(); err != nil {
// 		log.Fatal(err)
// 	}

// 	testFile := data.createFile("error_test.go.tmpl")
// 	data.createTemplate("error_test.go.tmpl", testFile)
// 	if err := testFile.Close(); err != nil {
// 		log.Fatal(err)
// 	}
// }
