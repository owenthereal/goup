package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show installed Go",
	Long:  "Show installed Go versions.",
	RunE:  runShow,
}

func runShow(cmd *cobra.Command, args []string) error {
	vers, err := listGoVers()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Version", "Active"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	for _, ver := range vers {
		if ver.Current {
			table.Append([]string{ver.Ver, "*"})
		} else {
			table.Append([]string{ver.Ver, ""})
		}
	}

	table.Render()

	return nil
}

type goVer struct {
	Ver     string
	Current bool
}

func listGoVers() ([]goVer, error) {
	files, err := ioutil.ReadDir(goupDir())
	if err != nil {
		return nil, err
	}

	current, err := currentGoVersion()
	if err != nil {
		return nil, err
	}

	var vers []goVer
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "go") {
			vers = append(vers, goVer{
				Ver:     strings.TrimPrefix(file.Name(), "go"),
				Current: current == file.Name(),
			})
		}
	}

	return vers, nil
}

func currentGoVersion() (string, error) {
	current := currentGoRootDir()
	goroot, err := os.Readlink(current)
	if err != nil {
		return "", err
	}

	return filepath.Base(goroot), nil
}
