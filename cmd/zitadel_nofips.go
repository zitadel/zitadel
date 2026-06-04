//go:build !fips

package cmd

import "github.com/spf13/viper"

func mergeFipsDefaultConfig(*viper.Viper) error { return nil }
