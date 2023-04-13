package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tnierman/backplane-tools/cmd/install"
)

var cmd = cobra.Command{
	Use: "backplane-tool",
	Short: "An OpenShift tool manager",
	Long: "This applications manages the tools needed to interact with OpenShift clusters",
	RunE: help,
}

func help(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

// Add subcommands
func init() {
	cmd.AddCommand(install.Cmd())
}

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
