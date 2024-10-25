package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	defaultCLIDocsRootDir = "./docs/cli"
)

var DocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Documentation for k8s-tools",
}

var docsGenerateCommand = &cobra.Command{
	Use:   "generate [output dir]",
	Short: "Generate documentation",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir := defaultCLIDocsRootDir
		if len(args) > 0 {
			outputDir = args[0]
		}

		err := generateMarkdownDocs(outputDir)
		if err != nil {
			log.Errorf("Failed to generate documentation: %v", err)
		}
		return err
	},
}

func generateMarkdownDocs(outputDir string) error {
	log.Infof("Generating markdown docs for CLI in directory: %v", outputDir)
	err := os.MkdirAll(outputDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	return doc.GenMarkdownTree(rootCmd, outputDir)
}

func init() {
	DocsCmd.AddCommand(docsGenerateCommand)
	rootCmd.DisableAutoGenTag = true
}

func init() {
	rootCmd.DisableAutoGenTag = true
	rootCmd.AddCommand(DocsCmd)
}
