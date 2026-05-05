// Package cmd defines the root command and shared
// functionality for the gh-readme CLI tool.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/givensuman/gh-readme/api"
	"github.com/spf13/cobra"
)

var ref string

var rootCmd = &cobra.Command{
	Use:   "gh-readme <owner/repo>",
	Short: "Render a GitHub repository's README in the terminal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("argument must be in owner/repo form, got %q", args[0])
		}

		owner, repo := parts[0], parts[1]

		content, err := api.GetReadme(owner, repo, ref)
		if err != nil {
			return err
		}

		rendered, err := glamour.Render(content, "auto")
		if err != nil {
			return fmt.Errorf("rendering README: %w", err)
		}

		return page(rendered)
	},
}

// page pipes text through $PAGER.
func page(text string) error {
	pager, isSet := os.LookupEnv("PAGER")
	if !isSet || pager == "" {
		pager = "less"
	}

	parts := strings.Fields(pager)
	name := parts[0]
	args := parts[1:]

	// if pager is "less" and -R isn't already set, add it so ANSI colours show
	if name == "less" {
		hasR := false
		for _, a := range args {
			if strings.Contains(a, "R") {
				hasR = true
				break
			}
		}

		if !hasR {
			args = append(args, "-R")
		}
	}

	c := exec.Command(name, args...)
	c.Stdin = strings.NewReader(text)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		// fall back to plain stdout
		fmt.Print(text)
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&ref, "ref", "", "Branch, tag, or commit SHA to fetch README from")
}
