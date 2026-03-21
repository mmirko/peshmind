/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the Prolog engine",
	Long: `Spawn a Prolog engine and serve it over a TCP socket. peshmind 
	supervises the execution of the engine until it terminates. `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := c.SpawnPrologEngine(kbPool); err != nil {
			cmd.Println("Failed to spawn Prolog engine:")
			cmd.Println(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
