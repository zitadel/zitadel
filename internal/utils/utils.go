package utils

import (
	"regexp"
	"strings"
)

func SanitizeDomain(name, suffix string) string {
	// References: https://www.nic.ad.jp/timeline/en/20th/appendix1.html
	// References:

	// - Zitadel right now replaces spaces in org name with hyphens
	label := strings.ReplaceAll(name, " ", "-")

	// - The label must be sanitized so it only contains alphanumeric characters and hyphens
	// - Invalid characters are replaced with and empty space
	label = string(regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAll([]byte(label), []byte("")))

	// - The label cannot exceed 63 characters
	if len(label) > 63 {
		label = label[:63]
	}

	// - The total length of the resulting domain can't exceed 253 characters
	domain := label + "." + suffix
	if len(domain) > 253 {
		truncateNChars := len(domain) - 253
		label = label[:len(label)-truncateNChars]
	}

	// - A domain label can't start with a hyphen
	if len(label) > 0 && label[0:1] == "-" {
		label = label[1:]
	}

	// - A domain label can't end with a hyphen
	if len(label) > 0 && label[len(label)-1:] == "-" {
		label = label[:len(label)-1]
	}

	domain = label + "." + suffix

	return strings.ToLower(domain)
}
