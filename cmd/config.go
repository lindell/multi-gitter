package cmd

import (
	"fmt"

	"github.com/lindell/multi-gitter/cmd/namedflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func configureConfig(flags namedflag.Set) {
	flags.StringP("config", "", "", "Path of the config file.")
}

func initializeConfig(cmd *cobra.Command) error {
	// Prioritize reading config files defined with --config
	if err := initializeDynamicConfig(cmd); err != nil {
		return err
	}

	// Read any config defined in static config files
	return initializeStaticConfig(cmd)
}

func initializeDynamicConfig(cmd *cobra.Command) error {
	configFile, _ := cmd.Flags().GetString("config")
	if configFile == "" {
		return nil
	}

	v := viper.New()

	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	bindFlags(cmd, v)

	return nil
}

func initializeStaticConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath("$HOME/.multi-gitter")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)

			switch val := val.(type) {
			case []interface{}:
				for _, v := range val {
					_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", v))
				}
			default:
				_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			}
		}
	})
}
