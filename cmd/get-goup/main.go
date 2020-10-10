package main

import (
	"github.com/owenthereal/goup/internal/commands"
	"github.com/sirupsen/logrus"
)

func main() {
	rootCmd := commands.NewInstallerCommand()
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
