package domain

import (
	"errors"
	"golang.org/x/text/language"

	z_errors "github.com/zitadel/zitadel/internal/errors"
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
	return !languagesAreContained(left, right)
}

func LanguageIsAllowed(allowUndefined bool, allowedLanguages []language.Tag, lang language.Tag) error {
	err := LanguageIsDefined(lang)
	if err != nil && allowUndefined {
		return nil
	}
	if err != nil {
		return err
	}
	if len(allowedLanguages) > 0 && !languageIsContained(allowedLanguages, lang) {
		return z_errors.ThrowPreconditionFailed(nil, "LANG-2M9fs", "Errors.Language.NotAllowed")
	}
	return nil
}

func LanguagesAreSupported(supportedLanguages []language.Tag, lang ...language.Tag) error {
	unsupported := make([]language.Tag, 0)
	for _, l := range lang {
		if l.IsRoot() {
			continue
		}
		if !languageIsContained(supportedLanguages, l) {
			unsupported = append(unsupported, l)
		}
	}
	if len(unsupported) == 0 {
		return nil
	}
	if len(unsupported) == 1 {
		return z_errors.ThrowInvalidArgument(nil, "LANG-lg4DP", "Errors.Language.NotSupported")
	}
	return z_errors.ThrowInvalidArgumentf(nil, "LANG-XHiK5", "Errors.Languages.NotSupported: %s", LanguagesToStrings(unsupported))
}

func LanguageIsDefined(lang language.Tag) error {
	if lang.IsRoot() {
		return z_errors.ThrowInvalidArgument(nil, "LANG-3M9f2", "Errors.Language.Undefined")
	}
	return nil
}

// LanguagesHaveDuplicates returns an error if the passed slices contains duplicates.
// The error lists the duplicates.
func LanguagesHaveDuplicates(langs []language.Tag) error {
	unique := make(map[language.Tag]struct{})
	duplicates := make([]language.Tag, 0)
	for _, lang := range langs {
		if _, ok := unique[lang]; ok {
			duplicates = append(duplicates, lang)
		}
		unique[lang] = struct{}{}
	}
	if len(duplicates) == 0 {
		return nil
	}
	if len(duplicates) > 1 {
		return z_errors.ThrowInvalidArgument(nil, "LANG-3M9f2", "Errors.Language.Duplicate")
	}
	return z_errors.ThrowInvalidArgumentf(nil, "LANG-XHiK5", "Errors.Languages.Duplicate: %s", LanguagesToStrings(duplicates))
}

func ParseLanguage(lang ...string) (tags []language.Tag, err error) {
	for _, l := range lang {
		t, parseErr := language.Parse(l)
		err = errors.Join(err, parseErr)
		tags = append(tags, t)
	}
	if err != nil {
		err = z_errors.ThrowInvalidArgument(err, "LANG-jc8Sq", "Errors.Language.NotParsed")
	}
	return tags, err
}

func languagesAreContained(languages, search []language.Tag) bool {
	for _, s := range search {
		if !languageIsContained(languages, s) {
			return false
		}
	}
	return true
}

func languageIsContained(languages []language.Tag, search language.Tag) bool {
	for _, lang := range languages {
		if lang == search {
			return true
		}
	}
	return false
}
