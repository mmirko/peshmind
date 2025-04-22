/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the switches in the database",
	Long:  `List the switches in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(c.Switches) != 0 {
			fmt.Println("Switches Database:")
			for switchName, sw := range c.Switches {
				fmt.Printf("\t%s\n", switchName)
				if sw.Description != "" {
					fmt.Printf("\t\tDescription: %s\n", sw.Description)
				}
				fmt.Printf("\t\tIP: %s\n", sw.IP)
				fmt.Printf("\t\tUsername: %s\n", sw.Username)
				fmt.Printf("\t\tPassword: %s\n", sw.Password)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
