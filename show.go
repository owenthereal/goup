package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current version of installed Go",
	RunE:  runShow,
}

func runShow(cmd *cobra.Command, args []string) error {
	current, err := currentGoRoot()
	if err != nil {
		return err
	}

	goroot, err := os.Readlink(current)
	if err != nil {
		return err
	}

	fmt.Println(filepath.Base(goroot))

	return nil
}
