package text

import (
	"golang.org/x/text/language"
)

func LanguageTagsToStrings(langs []language.Tag) []string {
	result := make([]string, len(langs))
	for i, lang := range langs {
		result[i] = lang.String()
	}
	return result
}
