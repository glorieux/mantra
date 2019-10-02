package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mantra",
		Short: "mantra - Let's build something",
		Long:  ``,
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of mantra",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("0.0.1")
		},
	})

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Creates new application",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createApplication(args[0])
		},
	}

	newCmd.AddCommand(&cobra.Command{
		Use:   "service",
		Short: "Create a new service",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createService(args[0])
		},
	})

	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
