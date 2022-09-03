package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <VERSION>...",
		Aliases: []string{"rm"},
		Short:   "Remove Go with a version",
		Long:    "Remove Go by providing a version.",
		Example: `
  goup remove 1.15.2
  goup remove 1.16.1 1.16.2
`,
		RunE: runRemove,
	}
}

func runRemove(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No version is specified")
	}

	for _, ver := range args {
		logger.Printf("Removing %s", ver)

		if !strings.HasPrefix(ver, "go") {
			ver = "go" + ver

		}

		if err := os.RemoveAll(GoupDir(ver)); err != nil {
			return err
		}
	}

	return nil
}
