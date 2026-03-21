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
		rawlist, _ := cmd.Flags().GetBool("rawlist")
		if len(c.Switches) != 0 {
			if rawlist {
				for switchName := range c.Switches {
					fmt.Print(switchName, " ")
				}
				fmt.Println()
			} else {
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
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("rawlist", "r", false, "List only the names of the switches in the database")
}
