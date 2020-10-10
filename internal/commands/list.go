package commands

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls-ver [regexp]",
		Short: `List Go versions to install`,
		Long: `List available Go versions matching a regexp filter for installation. If no filter is provided,
list all available versions.`,
		Example: `
  goup ls-ver
  goup ls-ver 1.15
`,
		RunE: runList,
	}
}

func runList(cmd *cobra.Command, args []string) error {
	var regexp string
	if len(args) > 0 {
		regexp = args[0]
	}

	vers, err := listGoVersions(regexp)
	if err != nil {
		return err
	}

	for _, ver := range vers {
		fmt.Println(ver)
	}

	return nil
}

func listGoVersions(re string) ([]string, error) {
	if re == "" {
		re = ".+"
	} else {
		re = fmt.Sprintf(`.*%s.*`, re)

	}

	cmd := exec.Command("git", "ls-remote", "--tags", "https://github.com/golang/go")
	refs, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(fmt.Sprintf(`refs/tags/go(%s)`, re))
	match := r.FindAllStringSubmatch(string(refs), -1)
	if match == nil {
		return nil, fmt.Errorf("No Go version found")
	}

	var vers []string
	for _, m := range match {
		vers = append(vers, m[1])
	}

	return vers, nil
}
