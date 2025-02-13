/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ConfigureCmd represents the config command
	ConfigureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Guides you through configuring Zitadel",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("config called")
			fmt.Println(viper.AllSettings())
			fmt.Println(viper.Sub("database").AllSettings())
			viper.en
		},
	}

	upgrade bool
)

func init() {
	// Here you will define your flags and configuration settings.
	ConfigureCmd.Flags().BoolVarP(&upgrade, "upgrade", "u", false, "Only changed configuration values since the previously used version will be asked for")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
