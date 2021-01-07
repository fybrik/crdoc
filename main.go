package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pkg "github.com/roee88/crdoc/pkg/builder"
)

const (
	outputOption    = "output"
	templateOption  = "template"
	resourcesOption = "resources"
	tocOption       = "toc"
)

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "crdoc",
		Short:         "Output markdown documentation from Kubernetes CustomResourceDefinition YAML files",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			outputOptionValue := viper.GetString(outputOption)
			templateOptionValue := viper.GetString(templateOption)
			resourcesOptionValue := viper.GetString(resourcesOption)
			tocOptionValue := viper.GetString(tocOption)

			model, err := pkg.LoadModel(tocOptionValue)
			if err != nil {
				return err
			}

			builder := pkg.ModelBuilder{
				Model:              model,
				Strict:             tocOptionValue != "",
				TemplatesDirOrFile: templateOptionValue,
				OutputFilepath:     outputOptionValue,
			}

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

			err = builder.Write()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP(outputOption, "o", "", "Path to output markdown file (required)")
	_ = cmd.MarkFlagRequired(outputOption)
	cmd.Flags().StringP(templateOption, "t", "", "Path to template file or directory (required)")
	_ = cmd.MarkFlagRequired(templateOption)
	cmd.Flags().StringP(resourcesOption, "r", "", "Path to directory with CustomResourceDefinition YAML files (required)")
	_ = cmd.MarkFlagRequired(resourcesOption)
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
