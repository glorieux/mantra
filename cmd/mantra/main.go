package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"glorieux.io/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mantra",
		Short: "mantra is a service based framework",
		Long:  ``,
	}

	v, _ := version.New("0.0.1")
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of mantra",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(v)
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
