package commands

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
)

func searchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [REGEXP]",
		Short: `Search Go versions to install`,
		Long: `Search available Go versions matching a regexp filter for installation. If no filter is provided,
list all available versions.`,
		Example: `
  goup search
  goup search 1.15
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

	cmd := exec.Command("git", "ls-remote", "--sort=version:refname", "--tags", "https://github.com/golang/go")
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
