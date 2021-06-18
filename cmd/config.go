package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func configureConfig(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "", "", "Path of the config file.")
}

func initializeConfig(cmd *cobra.Command) error {
	configFile, _ := cmd.Flags().GetString("config")
	if configFile == "" {
		return nil
	}

	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigFile(configFile)

	v.SetConfigType("yaml")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Bind the current command's flags to viper
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
