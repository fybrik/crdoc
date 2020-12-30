package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pkg "github.com/roee88/crdoc/pkg/loader"
)

const (
	outputDirOption    = "output-dir"
	configDirOption    = "config-dir"
	resourcesDirOption = "resources-dir"
	tocFileOption      = "toc-file"
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
			outputDir := viper.GetString(outputDirOption)
			configDir := viper.GetString(configDirOption)
			resourcesDir := viper.GetString(resourcesDirOption)
			tocFile := viper.GetString(tocFileOption)

			model, err := pkg.LoadModel(tocFile)
			if err != nil {
				return err
			}

			builder := pkg.ModelBuilder{
				Model:        model,
				Strict:       tocFile != "",
				TemplatesDir: path.Join(configDir, "templates"),
				OutputDir:    outputDir,
			}

			crds, err := pkg.LoadCRDs(resourcesDir)
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

	cmd.Flags().StringP(outputDirOption, "o", "", "Output directory for markdown file (required)")
	_ = cmd.MarkFlagRequired(outputDirOption)
	cmd.Flags().StringP(configDirOption, "c", "", "Configuration directory (required)")
	_ = cmd.MarkFlagRequired(configDirOption)
	cmd.Flags().StringP(resourcesDirOption, "r", "", "Directory containing CustomResourceDefinition YAML files (required)")
	_ = cmd.MarkFlagRequired(resourcesDirOption)
	cmd.Flags().StringP(tocFileOption, "t", "", "Path to table of contents YAML file")

	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

// Run the cli
func Run() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
