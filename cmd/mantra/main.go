package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"pkg.glorieux.io/mantra"
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
			fmt.Printf("v%s\n", mantra.VERSION)
		},
	})

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Creates new things",
		Args:  cobra.MinimumNArgs(1),
	}

	newCmd.AddCommand(&cobra.Command{
		Use:   "app",
		Short: "Create a new application",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Print("Please input the application's name: ")
				var appName string
				_, err := fmt.Scanln(&appName)
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				return createApplication(appName)
			}
			return createApplication(args[0])
		},
	})

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
