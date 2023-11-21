package domain

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
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
	if left == nil && right == nil {
		return false
	}
	if left == nil || right == nil || len(left) != len(right) {
		return true
	}
	return len(FilterOutLanguages(left, right)) > 0 || len(FilterOutLanguages(right, left)) > 0
}

// FilterOutLanguages returns a new slice of languages without the languages to exclude.
func FilterOutLanguages(originalLanguages, excludeLanguages []language.Tag) []language.Tag {
	filteredLanguages := make([]language.Tag, 0, len(originalLanguages))
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

func LanguageIsAllowed(allowUndefined bool, allowedLanguages []language.Tag, lang language.Tag) error {
	if err := LanguageIsDefined(lang); err != nil && !allowUndefined {
		return err
	}
	if len(allowedLanguages) > 0 && !languageIsContained(allowedLanguages, lang) {
		return errors.ThrowPreconditionFailed(nil, "LANG-2M9fs", "Errors.Language.NotAllowed")
	}
	return nil
}

func LanguagesAreSupported(lang ...language.Tag) error {
	unsupported := make([]language.Tag, 0)
	for _, l := range lang {
		if l.IsRoot() {
			continue
		}
		if !languageIsContained(i18n.SupportedLanguages(), l) {
			unsupported = append(unsupported, l)
		}
	}
	if len(unsupported) == 0 {
		return nil
	}
	if len(unsupported) == 1 {
		return errors.ThrowInvalidArgument(nil, "LANG-lg4DP", "Errors.Language.NotSupported")
	}
	return errors.ThrowInvalidArgumentf(nil, "LANG-XHiK5", "Errors.Languages.NotSupported: %s", LanguagesToStrings(unsupported))
}

func LanguageIsDefined(lang language.Tag) error {
	if lang.IsRoot() {
		return errors.ThrowInvalidArgument(nil, "LANG-3M9f2", "Errors.Language.Undefined")
	}
	return nil
}

func languageIsContained(languages []language.Tag, search language.Tag) bool {
	for _, lang := range languages {
		if lang == search {
			return true
		}
	}
	return false
}
