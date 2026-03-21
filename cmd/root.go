/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mmirko/peshmind/pkg/peshmind"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "peshmind",
	Short: "Prolog-based Engine for System & Host Management and INference Daemon",
	Long:  `Peshmind is a Go-based application that uses Prolog for system and host management.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cfgFile != "" {
			err := c.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("error loading config file: %w", err)
			}
		}
		return nil
	},
}

var cfgFile string // config file
var kbPool string  // Knowledge base pool, used to store the knowledge base in form of a Prolog files

var c *peshmind.Config // Config object

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "peshmind.json", "config file (default peshmind.json)")

	rootCmd.PersistentFlags().StringVarP(&kbPool, "kbpool", "k", "kbpool", "Knowledge base pool directory (default kbpool)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	c = peshmind.NewConfig()
}
