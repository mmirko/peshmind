/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate the data of a group of switches",
	Long:  `Simulate the data of a group of switches`,
	Run: func(cmd *cobra.Command, args []string) {
		if simName == "" {
			fmt.Println("Simulation name is required")
			return
		}
		if err := c.Simulate(simName); err != nil {
			fmt.Println("Error:", err)
		}
	},
}

var simName string // Name of the simulation to run

func init() {
	rootCmd.AddCommand(simulateCmd)

	simulateCmd.Flags().StringVarP(&simName, "simname", "s", "", "Name of the simulation to run")
	simulateCmd.MarkFlagRequired("simname")
}
