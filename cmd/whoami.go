/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"

	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show the current user info",
	Long:  `Show the current user info`,
	Run: func(cmd *cobra.Command, args []string) {
		u, err := requests.GetUser()
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println(u.Name)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoamiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoamiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
