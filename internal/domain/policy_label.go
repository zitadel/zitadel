package domain

import (
	"regexp"
	"unicode/utf8"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var colorRegex = regexp.MustCompile("^$|^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")

type LabelPolicy struct {
	models.ObjectRoot

	State   LabelPolicyState
	Default bool

	PrimaryColor    string
	BackgroundColor string
	WarnColor       string
	FontColor       string
	LogoURL         string
	IconURL         string

	PrimaryColorDark    string
	BackgroundColorDark string
	WarnColorDark       string
	FontColorDark       string
	LogoDarkURL         string
	IconDarkURL         string

	Font string

	HideLoginNameSuffix bool
	ErrorMsgPopup       bool
	DisableWatermark    bool
	ThemeMode           LabelPolicyThemeMode

	LinkPreviewTitle       string
	LinkPreviewDescription string
	LinkPreviewImageURL    string
}

type LabelPolicyState int32

const (
	LabelPolicyStateUnspecified LabelPolicyState = iota
	LabelPolicyStateActive
	LabelPolicyStateRemoved
	LabelPolicyStatePreview

	labelPolicyStateCount
)

type LabelPolicyThemeMode int32

const (
	LabelPolicyThemeAuto LabelPolicyThemeMode = iota
	LabelPolicyThemeLight
	LabelPolicyThemeDark
)

func (f LabelPolicy) IsValid() error {
	if !colorRegex.MatchString(f.PrimaryColor) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-391dG", "Errors.Policy.Label.Invalid.PrimaryColor")
	}
	if !colorRegex.MatchString(f.BackgroundColor) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-502F1", "Errors.Policy.Label.Invalid.BackgroundColor")
	}
	if !colorRegex.MatchString(f.WarnColor) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-nvw33", "Errors.Policy.Label.Invalid.WarnColor")
	}
	if !colorRegex.MatchString(f.FontColor) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-93mSf", "Errors.Policy.Label.Invalid.FontColor")
	}
	if !colorRegex.MatchString(f.PrimaryColorDark) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-391dG", "Errors.Policy.Label.Invalid.PrimaryColorDark")
	}
	if !colorRegex.MatchString(f.BackgroundColorDark) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-llsp2", "Errors.Policy.Label.Invalid.BackgroundColorDark")
	}
	if !colorRegex.MatchString(f.WarnColorDark) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-2b6sf", "Errors.Policy.Label.Invalid.WarnColorDark")
	}
	if !colorRegex.MatchString(f.FontColorDark) {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-3M0fs", "Errors.Policy.Label.Invalid.FontColorDark")
	}
	return ValidateLinkPreviewText(f.LinkPreviewTitle, f.LinkPreviewDescription)
}

// ValidateLinkPreviewText checks the length of the link preview title and description in characters.
func ValidateLinkPreviewText(title, description string) error {
	if utf8.RuneCountInString(title) > 120 {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-ogT1l", "Errors.Policy.Label.Invalid.LinkPreviewTitle")
	}
	if utf8.RuneCountInString(description) > 300 {
		return zerrors.ThrowInvalidArgument(nil, "POLICY-ogD3s", "Errors.Policy.Label.Invalid.LinkPreviewDescription")
	}
	return nil
}
