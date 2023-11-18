package domain

import (
	"golang.org/x/text/language"
)

func StringsToLanguages(langs []string) []language.Tag {
	return GenericMapSlice(langs, language.Make)
}

func LanguagesToStrings(langs []language.Tag) []string {
	return GenericMapSlice(langs, func(lang language.Tag) string { return lang.String() })
}

func GenericMapSlice[T any, U any](from []T, mapTo func(T) U) []U {
	if from == nil {
		return nil
	}
	result := make([]U, len(from))
	for i, lang := range from {
		result[i] = mapTo(lang)
	}
	return result
}

// LanguagesDiffer returns true if the languages differ.
func LanguagesDiffer(left, right []language.Tag) bool {
	if len(left) != len(right) {
		return true
	}
	return len(FilterOutLanguages(left, right)) > 0 || len(FilterOutLanguages(right, left)) > 0
}

// FilterOutLanguages returns a new slice of languages without the languages to exclude.
func FilterOutLanguages(originalLanguages, excludeLanguages []language.Tag) []language.Tag {
	var filteredLanguages []language.Tag
	for _, lang := range originalLanguages {
		var found bool
		for _, excludeLang := range excludeLanguages {
			if lang == excludeLang {
				found = true
				break
			}
		}
		if !found {
			filteredLanguages = append(filteredLanguages, lang)
		}
	}
	return filteredLanguages
}
