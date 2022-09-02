package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.4.1"

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show goup version",
		RunE: func(c *cobra.Command, args []string) error {
			_, err := fmt.Printf("goup version v%s\n", Version)
			return err
		},
	}
}
