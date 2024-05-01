package cmd

import "github.com/spf13/pflag"

// stringSlice is a wrapped around *pflag.FlagSet.GetStringSlice to allow nil when the flag is not set
func stringSlice(set *pflag.FlagSet, name string) ([]string, error) {
	if !set.Changed(name) {
		return nil, nil
	}
	return set.GetStringSlice(name)
}
