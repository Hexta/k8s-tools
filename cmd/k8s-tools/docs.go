package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	cliDocsRootDir = "./docs/cli"
)

func newDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Documentation for k8s-tools",
	}

	cmd.AddCommand(newDocsGenerateCmd())

	return cmd
}

func newDocsGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := generateMarkdownDocs()
			if err != nil {
				log.Errorf("Failed to generate documentation: %v", err)
			}
			return err
		},
	}
}

func generateMarkdownDocs() error {
	log.Infof("Generating markdown docs for CLI in directory: %v", cliDocsRootDir)
	err := os.MkdirAll(cliDocsRootDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	return doc.GenMarkdownTree(rootCmd, cliDocsRootDir)
}

func init() {
	rootCmd.AddCommand(newDocsCmd())
}
