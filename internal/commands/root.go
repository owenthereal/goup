package commands

import (
	"fmt"
	"github.com/owenthereal/goup/internal/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thewolfnl/go-multiplechoice"
	"os"
	"path/filepath"
)

var (
	homedir string
	logger  *logrus.Logger

	ProfileFiles []string
)

func init() {
	logger = logrus.New()

	var err error
	homedir, err = os.UserHomeDir()
	if err != nil {
		logger.Fatal(err)
	}

	ProfileFiles = []string{
		filepath.Join(homedir, ".profile"),
		filepath.Join(homedir, ".zprofile"),
		filepath.Join(homedir, ".bash_profile"),
	}
}

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "goup",
		Short: "The Go installer",
		RunE:  choiceVersion,
	}

	rootCmd.AddCommand(installCmd())
	rootCmd.AddCommand(removeCmd())
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(showCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(versionCmd())

	return rootCmd
}

func GoupBinDir() string {
	return GoupDir("bin")
}

func GoupCurrentDir() string {
	return GoupDir("current")
}

func GoupEnvFile() string {
	return GoupDir("env")
}

func GoupCurrentBinDir() string {
	return GoupDir("current", "bin")
}

func goupVersionDir(ver string) string {
	return GoupDir(ver)
}

func GoupDir(paths ...string) string {
	elem := []string{homedir, ".go"}
	elem = append(elem, paths...)

	return filepath.Join(elem...)
}

func choiceVersion(cmd *cobra.Command, args []string) error {
	vers, err := listGoVers()
	if err != nil {
		return err
	}

	if len(vers) == 0 {
		showGoIfExist()
		return nil
	}

	var curVer string

	var options = make([]string, 0, len(vers))
	for _, v := range vers {
		options = append(options, v.Ver)
		if v.Current {
			curVer = v.Ver
		}
	}

	// choice version
	ver := MultipleChoice.Selection(fmt.Sprintf("CurVer %s \nUse up/down arrow keys to select a version, return key to install \n", color.Str2Cyan(curVer)), options[:])

	fmt.Println()
	if ver == curVer {
		fmt.Println(color.Str2Red(fmt.Sprintf("installed：go v%s \n", ver)))
		return nil
	}

	if err := symlink("go" + ver); err != nil {
		return err
	}
	fmt.Println(color.Str2Red(fmt.Sprintf("installed：go v%s \n", ver)))
	return nil
}
