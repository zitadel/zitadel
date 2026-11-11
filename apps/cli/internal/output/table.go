package output

import (
	"os"

	"github.com/pterm/pterm"
)

// Table writes rows as a formatted table to stdout.
func Table(header []string, rows [][]string) {
	if len(rows) == 0 {
		pterm.Info.Println("No results found.")
		return
	}

	tableData := pterm.TableData{header}
	for _, row := range rows {
		tableData = append(tableData, row)
	}

	// Disable styling if stdout is piped to keep it raw for scripts
	if IsStdoutPiped() {
		pterm.DisableStyling()
	}

	pterm.DefaultTable.
		WithHasHeader().
		WithData(tableData).
		WithWriter(os.Stdout).
		Render()
}
