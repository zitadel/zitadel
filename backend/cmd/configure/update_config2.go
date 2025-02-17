package configure

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/manifoldco/promptui"
)

func Update2(config any) {
	// Print the intro
	printIntro()

	// Start the interactive CLI
	interactiveCLI(reflect.ValueOf(config), 0)
	fmt.Println(config)
}

const (
	ExitValue = "<Exit>"
	BackValue = "⬅ Back"
	Prefix    = "📁 "
)

var introTemplate = `
   +----------------------------------------+
   |      🛠  Config Interactive CLI  🛠      |
   +----------------------------------------+
   |                                        |
   | %5s : Dive into nested config       |
   | %6s : Return to previous menu       |
   | %6s : Exit application              |
   |                                        |
   |   Choose an option to explore!         |
   |                                        |
   +----------------------------------------+
   `

func printIntro() {
	fmt.Printf(introTemplate, Prefix, BackValue, ExitValue)
}

// interactiveCLI handles the interactive CLI
func interactiveCLI(v reflect.Value, depth int) {
	for {
		var items []string

		// If depth is greater than 0, we are in a nested struct and should add a "Back" option
		if depth > 0 {
			items = append(items, BackValue)
		}

		// Add all the field names
		items = append(items, getFieldNames(v)...)

		// Add an "Exit" option
		items = append(items, ExitValue)

		prompt := promptui.Select{
			Label: "Select Field",
			Items: items,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case BackValue:
			return
		case ExitValue:
			// Exit the entire application
			os.Exit(0)
		default:
			fieldName := strings.TrimPrefix(result, Prefix)
			selectedField := v.FieldByName(fieldName)
			if selectedField.Kind() == reflect.Struct {
				interactiveCLI(selectedField, depth+1)
			} else {
				prompt := promptui.Prompt{
					Label:   fmt.Sprintf("Field %s (%s)", result, selectedField.Kind()),
					Default: fmt.Sprintf("%v", selectedField.Interface()),
				}
				res, err := prompt.Run()
				fmt.Println(res, err)
				// fmt.Printf("%s: %v\n", result, selectedField.Interface())
			}
		}
	}
}

// getFieldNames returns all the field names
func getFieldNames(v reflect.Value) []string {
	t := v.Type()
	var fieldNames []string
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		if v.Field(i).Kind() == reflect.Struct {
			fieldName = Prefix + fieldName
		}
		fieldNames = append(fieldNames, fieldName)
	}
	return fieldNames
}
