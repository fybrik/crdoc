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
	strictOption    = "strict"
	fileOption      = "file"
	configDirOption = "config-dir"
	outputDirOption = "output-dir"
)

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "crdoc",
		Short:         "Output specification for CustomResourceDefinition",
		Long:          "Output the CustomResourceDefinition specification in a format usable for the Kubernetes website",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			confDir := viper.GetString(configDirOption)
			outDir := viper.GetString(outputDirOption)
			crdFile := viper.GetString(fileOption)
			strictMode := viper.GetBool(strictOption)

			model, err := pkg.LoadModel(path.Join(confDir, "toc.yaml"))
			if err != nil {
				return err
			}

			builder := pkg.ModelBuilder{
				Model:        model,
				Strict:       strictMode,
				TemplatesDir: path.Join(confDir, "templates"),
				OutputDir:    outDir,
			}

			crd, err := pkg.LoadCRD(crdFile)
			if err != nil {
				return err
			}
			err = builder.Add(crd)
			if err != nil {
				return err
			}

			err = builder.Write()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP(configDirOption, "c", "", "Directory containing documentation configuration")
	_ = cmd.MarkFlagRequired(configDirOption)
	cmd.Flags().StringP(outputDirOption, "o", "", "Directory to write markdown files")
	_ = cmd.MarkFlagRequired(outputDirOption)
	cmd.PersistentFlags().StringP(fileOption, "f", "", "CustomResourceDefinition YAML file")
	_ = cmd.MarkFlagRequired(fileOption)
	cmd.PersistentFlags().BoolP(strictOption, "s", false, "Strict mode")

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
