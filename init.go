package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	envFileContent = `export PATH="$HOME/.go/bin":"$HOME/.go/current/bin:$PATH`

	instruction = `To get started you need goup's ($HOME/.go/bin) and Go's bin directory ($HOME/.go/current/bin)
in your PATH environment variable. Add the following to your shell startup script:

source $HOME/.go/env

To configure your current shell run source $HOME/.go/env`
)

var (
	initUpdateFlag bool

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the goup environment file.",
		RunE:  runInit,
	}
)

func init() {
	initCmd.PersistentFlags().BoolVar(&initUpdateFlag, "update", false, "Update Go to the latest")
}

func runInit(cmd *cobra.Command, args []string) error {
	if initUpdateFlag {
		if err := runUpdate(cmd, args); err != nil {
			return err
		}
	}

	envFile, err := envFile()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(envFile), 0755); err != nil {
		return err
	}

	// ignore error, similar to rm -f
	os.Remove(envFile)

	if err := ioutil.WriteFile(envFile, []byte(envFileContent), 0664); err != nil {
		return err
	}

	fmt.Println(instruction)

	return nil
}

func envFile() (string, error) {
	home, err := homedir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	return filepath.Join(home, ".go", "env"), nil
}
