package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Table writes rows as a tab-aligned table to stdout.
func Table(header []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(header, "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}
