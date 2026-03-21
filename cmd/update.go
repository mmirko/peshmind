/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the database of a switch",
	Long:  `Update the database of a switch`,
	Run: func(cmd *cobra.Command, args []string) {
		if switchName == "" {
			fmt.Println("Switch name or \"all\" is required")
			return
		}

		swToUpdate := make([]string, 0)

		if switchName != "all" {
			if _, ok := c.Switches[switchName]; !ok {
				fmt.Printf("Switch %s not found\n", switchName)
				return
			} else {
				swToUpdate = append(swToUpdate, switchName)
			}
		} else {
			for name := range c.Switches {
				swToUpdate = append(swToUpdate, name)
			}
		}

		fmt.Print("Enter admin password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("\nFailed to read password:")
			fmt.Println(err.Error())
			return
		}
		fmt.Println() // Add newline after password input

		// Convert to string if needed
		passwordStr := string(password)
		// fmt.Println("Password entered:", passwordStr)

		for _, switchName := range swToUpdate {

			sw := c.Switches[switchName]
			if c.Switches[switchName].Password == "ask" {
				sw.Password = passwordStr
			} else {
				sw.Password = c.Switches[switchName].Password
			}

			if err := sw.ApplyTemplate(); err != nil {
				fmt.Println("Failed to apply template:")
				fmt.Println(err.Error())
				return
			} else {
				fmt.Printf("Template applied successfully for switch %s\n", switchName)
			}

			if err := sw.FetchData(); err != nil {
				fmt.Println("Failed to fetch data:")
				fmt.Println(err.Error())
				return
			} else {
				fmt.Printf("Data fetched successfully for switch %s\n", switchName)
			}

			if err := sw.SaveKB(kbPool); err != nil {
				fmt.Println("Failed to save KB:")
				fmt.Println(err.Error())
				return
			} else {
				fmt.Printf("KB saved successfully for switch %s\n", switchName)
			}
		}
	},
}

var switchName string // Name of the switch to update

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&switchName, "switch", "s", "", "Name of the switch to update")
	updateCmd.MarkFlagRequired("switch")
}
