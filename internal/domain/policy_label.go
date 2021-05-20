package domain

import (
	"regexp"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
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
}

type LabelPolicyState int32

const (
	LabelPolicyStateUnspecified LabelPolicyState = iota
	LabelPolicyStateActive
	LabelPolicyStateRemoved
	LabelPolicyStatePreview

	labelPolicyStateCount
)

func (f LabelPolicy) IsValid() error {
	if !colorRegex.MatchString(f.PrimaryColor) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-391dG", "Errors.Policy.Label.Invalid.PrimaryColor")
	}
	if !colorRegex.MatchString(f.BackgroundColor) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-502F1", "Errors.Policy.Label.Invalid.BackgroundColor")
	}
	if !colorRegex.MatchString(f.WarnColor) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-nvw33", "Errors.Policy.Label.Invalid.WarnColor")
	}
	if !colorRegex.MatchString(f.FontColor) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-93mSf", "Errors.Policy.Label.Invalid.FontColor")
	}
	if !colorRegex.MatchString(f.PrimaryColorDark) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-391dG", "Errors.Policy.Label.Invalid.PrimaryColorDark")
	}
	if !colorRegex.MatchString(f.BackgroundColorDark) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-llsp2", "Errors.Policy.Label.Invalid.BackgroundColorDark")
	}
	if !colorRegex.MatchString(f.WarnColorDark) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-2b6sf", "Errors.Policy.Label.Invalid.WarnColorDark")
	}
	if !colorRegex.MatchString(f.FontColorDark) {
		return caos_errs.ThrowInvalidArgument(nil, "POLICY-3M0fs", "Errors.Policy.Label.Invalid.FontColorDark")
	}
	return nil
}

func (f LabelPolicyState) Valid() bool {
	return f >= 0 && f < labelPolicyStateCount
}

func (s LabelPolicyState) Exists() bool {
	return s != LabelPolicyStateUnspecified && s != LabelPolicyStateRemoved
}
