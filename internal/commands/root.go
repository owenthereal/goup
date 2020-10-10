package commands

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	homedir string
	logger  *logrus.Logger
)

func init() {
	logger = logrus.New()

	var err error
	homedir, err = os.UserHomeDir()
	if err != nil {
		logger.Fatal(err)
	}
}

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "goup",
		Short: "The Go installer",
	}

	rootCmd.AddCommand(installCmd())
	rootCmd.AddCommand(removeCmd())
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(showCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(versionCmd())

	return rootCmd
}

func goupBinDir() string {
	return goupDir("bin")
}

func currentGoRootDir() string {
	return goupDir("current")
}

func currentGoBinDir() string {
	return goupDir("current", "bin")
}

func versionGoRootDir(ver string) string {
	return goupDir(ver)
}

func envFile() string {
	return goupDir("env")
}

func goupDir(paths ...string) string {
	elem := []string{homedir, ".go"}
	elem = append(elem, paths...)

	return filepath.Join(elem...)
}
