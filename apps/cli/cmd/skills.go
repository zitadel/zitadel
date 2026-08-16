package cmd

import (
	"fmt"
	"io/fs"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/skills"
)

func newSkillsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skills [topic]",
		Short: "Display agent-readable context and skill documentation",
		Long: `Dump structured Markdown skill files for AI agents.

Without arguments, displays the general CONTEXT.md file.
With a topic argument, displays the skill file for that topic if available.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return printSkill("CONTEXT.md")
			}
			return printSkill(args[0] + ".md")
		},
	}
}

func printSkill(filename string) error {
	data, err := fs.ReadFile(skills.FS, filename)
	if err != nil {
		entries, _ := fs.ReadDir(skills.FS, ".")
		var names []string
		for _, e := range entries {
			if !e.IsDir() && e.Name() != "embed.go" {
				names = append(names, e.Name())
			}
		}
		return fmt.Errorf("skill %q not found. Available: %v", filename, names)
	}
	fmt.Print(string(data))
	return nil
}
