/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listSwitchesCmd represents the listSwitches command
var listSwitchesCmd = &cobra.Command{
	Use:   "listSwitches",
	Short: "List the switches in the database",
	Long:  `List the switches in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(c.Switches) != 0 {
			fmt.Println("Switches Database:")
			for switchName, ip := range c.Switches {
				fmt.Printf("\t%s, IP: %s\n", switchName, ip)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listSwitchesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listSwitchesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listSwitchesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
