package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <version>",
	Short: "Remove Go with a version",
	Long:  "Remove Go by providing a version.",
	RunE:  runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No version is specified")
	}

	ver := args[0]
	if !strings.HasPrefix(ver, "go") {
		ver = "go" + ver

	}

	return os.RemoveAll(goupDir(ver))
}
