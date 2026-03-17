/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mmirko/peshmind/pkg/peshmind"
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
		s := peshmind.NewSimulation(c)

		if err := s.Simulate(simName); err != nil {
			fmt.Println("Error:", err)
		}

		if emitDotFile != "" {
			if dotData, err := s.EmitDot(); err == nil {
				if err := os.WriteFile(emitDotFile, []byte(dotData), 0644); err != nil {
					fmt.Println("Error writing DOT file:", err)
				} else {
					if s.Debug {
						fmt.Println("DOT file written to", emitDotFile)
					}
				}
			} else {
				fmt.Println("Error emitting DOT:", err)
			}
		}

		if outputFile != "" {
			if outputData, err := s.EmitOutput(); err == nil {
				if err := os.WriteFile(outputFile, []byte(outputData), 0644); err != nil {
					fmt.Println("Error writing output file:", err)
				} else {
					if s.Debug {
						fmt.Println("Output file written to", outputFile)
					}
				}
			} else {
				fmt.Println("Error emitting output:", err)
			}
		}
	},
}

var simName string     // Name of the simulation to run
var emitDotFile string // Whether to emit a DOT file for visualization
var outputFile string  // Output file for the simulation results

func init() {
	rootCmd.AddCommand(simulateCmd)

	simulateCmd.Flags().StringVarP(&simName, "simname", "s", "", "Name of the simulation to run")
	simulateCmd.Flags().StringVarP(&emitDotFile, "emit-dot", "d", "", "Emit a DOT file for visualization")
	simulateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for the simulation results")
	simulateCmd.MarkFlagRequired("simname")
}
