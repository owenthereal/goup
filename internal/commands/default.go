package commands

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func defaultCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "default <VERSION>...",
		Short: "Set the default Go version",
		Long: `Set the default Go version to one specified. If no version is provided,
a prompt will show to select a installed Go version.`,
		Example: `
  goup default # A prompt will show to select a version
  goup default 1.15.2
`,
		RunE: runDefault,
	}
}

func runDefault(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return switchVer(args[0])
	}

	vers, err := listGoVers()
	if err != nil {
		return err
	}

	if len(vers) == 0 {
		showGoIfExist()
		return nil
	}

	var (
		pos int
	)

	var items = make([]string, 0, len(vers))

	for idx, v := range vers {
		items = append(items, v.Ver)
		if v.Current {
			pos = idx
		}
	}

	prompt := promptui.Select{
		Label:     "Select a version",
		Items:     items,
		CursorPos: pos,
	}

	_, ver, err := prompt.Run()
	if err != nil {
		return err
	}

	if err := switchVer(ver); err != nil {
		return err
	}

	return nil
}
