/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"
)

// UpgradeCmd represents the upgrade command
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades Zitadel from a previous version",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upgrade called")
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
