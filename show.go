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
	ver, err := currentGoVersion()
	if err != nil {
		return err
	}

	fmt.Println(ver)

	return nil
}

func currentGoVersion() (string, error) {
	current := currentGoRootDir()
	goroot, err := os.Readlink(current)
	if err != nil {
		return "", err
	}

	return filepath.Base(goroot), nil
}
