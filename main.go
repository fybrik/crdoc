// Copyright 2021 IBM Corp.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pkg "fybrik.io/crdoc/pkg/builder"
)

const (
	outputOption    = "output"
	templateOption  = "template"
	resourcesOption = "resources"
	tocOption       = "toc"
)

//go:embed templates/*
var builtinTemplates embed.FS

//go:embed VERSION
var version string

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crdoc",
		Short: "Output markdown documentation from Kubernetes CustomResourceDefinition YAML files",
		Example: `
  # Generate example/output.md from example/crds using the default markdown.tmpl template: 
  crdoc --resources example/crds --output example/output.md

  # Override template (builtin or custom):
  crdoc --resources example/crds --output example/output.md --template frontmatter.tmpl
  crdoc --resources example/crds --output example/output.md --template templates_folder/file.tmpl

  # Use a Table of Contents to filter and order CRDs
  crdoc --resources example/crds --output example/output.md --toc example/toc.yaml
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		Version: strings.TrimSpace(version),
		RunE: func(cmd *cobra.Command, args []string) error {
			outputOptionValue := viper.GetString(outputOption)
			templateOptionValue := viper.GetString(templateOption)
			resourcesOptionValue := viper.GetString(resourcesOption)
			tocOptionValue := viper.GetString(tocOption)

			model, err := pkg.LoadModel(tocOptionValue)
			if err != nil {
				return err
			}

			builder := pkg.NewModelBuilder(model, tocOptionValue != "",
				templateOptionValue, outputOptionValue, builtinTemplates)

			crds, err := pkg.LoadCRDs(resourcesOptionValue)
			if err != nil {
				return err
			}

			for _, crd := range crds {
				err = builder.Add(crd)
				if err != nil {
					return err
				}
			}

			err = builder.Output()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP(outputOption, "o", "", "Path to output markdown file (required)")
	_ = cmd.MarkFlagRequired(outputOption)
	cmd.Flags().StringP(resourcesOption, "r", "", "Path to YAML file or directory containing CustomResourceDefinitions (required)")
	_ = cmd.MarkFlagRequired(resourcesOption)
	cmd.Flags().StringP(templateOption, "t", "markdown.tmpl", "Path to file in a templates directory")
	cmd.Flags().StringP(tocOption, "c", "", "Path to table of contents YAML file")

	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	// Run the cli
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
