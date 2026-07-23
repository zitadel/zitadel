// Package skills provides embedded skill documentation for AI agents.
package skills

import "embed"

// FS contains the embedded skill files.
//
//go:embed *.md
var FS embed.FS
