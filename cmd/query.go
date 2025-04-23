/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/guregu/pengine"
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the Prolog engine",
	Long: `The query command allows you to send queries to the Prolog engine spawned by
the serve command. It will return the results of the query in a human-readable format.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := pengine.Client{
			URL:   c.EndPoint,
			Chunk: 1,
			Debug: false,
		}.Create(context.Background(), true)

		if err != nil {
			cmd.PrintErrf("Error creating client: %v\n", err)
			return
		}

		if len(args) > 0 {

			ctx := context.Background()

			for _, arg := range args {

				as, err := client.AskProlog(ctx, arg)
				if err != nil {
					cmd.PrintErrf("Error asking Prolog: %v\n", err)
					return
				}

				for as.Next(ctx) {
					fmt.Println(as.Current())
					if err := as.Err(); err != nil {
						cmd.PrintErrf("Error in answer set: %v\n", err)
						return
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
