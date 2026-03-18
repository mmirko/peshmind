/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/pengine"
	"github.com/mmirko/peshmind/pkg/peshmind"
	"github.com/spf13/cobra"
)

// dotCmd represents the dot command
var dotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Generate a dot file",
	Long: `The dot command generates a dot file from the Prolog engine.
	This file can be used to visualize several aspects of the network.
	For example, it can be used to visualize how the switches are connected.`,
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

		ctx := context.Background()

		as, err := client.Ask(ctx, "directp(X,Y,PORTX,PORTY)")
		if err != nil {
			cmd.PrintErrf("Error asking Prolog: %v\n", err)
			return
		}

		sws := make(map[string]*peshmind.SwitchDOT)
		links := make(map[string]struct{})

		for as.Next(ctx) {
			cur := as.Current()
			var x string
			if cur["X"].Atom != nil {
				x = *cur["X"].Atom
			} else if cur["X"].Number != nil {
				x = cur["X"].Number.String()
			}
			var y string
			if cur["Y"].Atom != nil {
				y = *cur["Y"].Atom
			} else if cur["Y"].Number != nil {
				y = cur["Y"].Number.String()
			}
			var portx string
			if cur["PORTX"].Atom != nil {
				portx = *cur["PORTX"].Atom
			} else if cur["PORTX"].Number != nil {
				portx = cur["PORTX"].Number.String()
			}
			var porty string
			if cur["PORTY"].Atom != nil {
				porty = *cur["PORTY"].Atom
			} else if cur["PORTY"].Number != nil {
				porty = cur["PORTY"].Number.String()
			}

			if _, ok := sws[x]; !ok {
				sws[x] = &peshmind.SwitchDOT{
					ID:    x,
					Ports: make(map[string]*peshmind.SwitchPort),
				}
			}
			if _, ok := sws[y]; !ok {
				sws[y] = &peshmind.SwitchDOT{
					ID:    y,
					Ports: make(map[string]*peshmind.SwitchPort),
				}
			}

			if _, ok := sws[x].Ports[portx]; !ok {
				sws[x].Ports[portx] = &peshmind.SwitchPort{
					SwitchID: x,
					Name:     portx,
					EndPoint: &peshmind.SwitchPort{
						SwitchID: y,
						Name:     porty,
					},
				}
				var link string
				if x < y {
					link = fmt.Sprintf("%s_%s -- %s_%s", san(x), portx, san(y), porty)
				} else {
					link = fmt.Sprintf("%s_%s -- %s_%s", san(y), porty, san(x), portx)
				}
				if _, ok := links[link]; !ok {
					links[link] = struct{}{}
				}
			}
			if _, ok := sws[y].Ports[porty]; !ok {
				sws[y].Ports[porty] = &peshmind.SwitchPort{
					SwitchID: y,
					Name:     porty,
					EndPoint: &peshmind.SwitchPort{
						SwitchID: x,
						Name:     portx,
					},
				}
				var link string
				if y < x {
					link = fmt.Sprintf("%s_%s -- %s_%s", san(y), porty, san(x), portx)
				} else {
					link = fmt.Sprintf("%s_%s -- %s_%s", san(x), portx, san(y), porty)
				}
				if _, ok := links[link]; !ok {
					links[link] = struct{}{}
				}
			}

			if err := as.Err(); err != nil {
				cmd.PrintErrf("Error in answer set: %v\n", err)
				return
			}
		}

		fmt.Println("graph G {")
		for _, sw := range sws {
			fmt.Printf("  subgraph cluster_%s {\n", san(sw.ID))
			fmt.Printf("    label=\"%s\"\n", sw.ID)
			for _, port := range sw.Ports {
				fmt.Printf("    %s_%s [label=\"%s\"];\n", san(sw.ID), port.Name, port.Name)
			}
			fmt.Println("  }")
		}
		for link := range links {
			fmt.Printf("  %s;\n", link)
		}
		fmt.Println("}")
	},
}

func san(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

func init() {
	rootCmd.AddCommand(dotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
