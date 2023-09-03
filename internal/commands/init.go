package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	GoupEnvFileContent       = `export PATH="$HOME/.go/bin:$HOME/.go/current/bin:$PATH"`
	ProfileFileSourceContent = `source "$HOME/.go/env"`

	welcomeTmpl = `Welcome to Goup!

Goup and Go will be located at:

  {{ .GoupDir }}

The Goup command will be located at:

  {{ .GoupBinDir }}

The go, gofmt and other Go commands will be located at:

  {{ .CurrentGoBinDir }}

To get started you need Goup's bin directory ({{ .GoupBinDir }}) and
Go's bin directory ({{ .CurrentGoBinDir }}) in your PATH environment
variable. These two paths will be added to your PATH environment variable by
modifying the profile files located at:
{{ range $index, $elem :=.ProfileFiles }}
  {{ $elem -}}
{{ end }}

Next time you log in this will be done automatically. To configure your
current shell run source $HOME/.go/env.
`
)

var (
	initCmdSkipInstallFlag bool
	initCmdSkipPromptFlag  bool
)

func initCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:    "init",
		Short:  "Initialize the goup environment file",
		Hidden: true,
		RunE:   runInit,
	}

	initCmd.PersistentFlags().BoolVar(&initCmdSkipInstallFlag, "skip-install", false, "Skip installing Go")
	initCmd.PersistentFlags().BoolVar(&initCmdSkipPromptFlag, "skip-prompt", false, "Skip confirmation prompt")

	return initCmd
}

func runInit(cmd *cobra.Command, args []string) error {
	tmpl, err := template.New("").Parse(welcomeTmpl)
	if err != nil {
		return err
	}

	params := struct {
		GoupDir         string
		GoupBinDir      string
		CurrentGoBinDir string
		ProfileFiles    []string
	}{
		GoupDir:         GoupDir(),
		GoupBinDir:      GoupBinDir(),
		CurrentGoBinDir: GoupCurrentBinDir(),
		ProfileFiles:    ProfileFiles,
	}
	if err := tmpl.Execute(os.Stdout, params); err != nil {
		return err
	}

	if !initCmdSkipPromptFlag {
		// Add a line break
		fmt.Println("")

		prompt := promptui.Prompt{
			Label:     "Would you like to proceed with the installation",
			IsConfirm: true,
		}

		if _, err := prompt.Run(); err != nil {
			return fmt.Errorf("interrupted")
		}

	}

	ef := GoupEnvFile()
	if err := os.MkdirAll(filepath.Dir(ef), 0755); err != nil {
		return err
	}

	// ignore error, similar to rm -f
	os.Remove(ef)

	if err := os.WriteFile(ef, []byte(GoupEnvFileContent), 0664); err != nil {
		return err
	}

	if err := appendSourceToProfiles(ProfileFiles); err != nil {
		return err
	}

	if !initCmdSkipInstallFlag {
		// Add a line break
		fmt.Println("")

		if err := runInstall(cmd, args); err != nil {
			return err
		}
	}

	return nil
}

func appendSourceToProfiles(profiles []string) error {
	for _, profile := range profiles {
		if err := appendToFile(profile, ProfileFileSourceContent); err != nil {
			return err
		}
	}

	return nil
}

func appendToFile(filename, value string) error {
	ok, err := checkStringExistsFile(filename, value)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + value + "\n"); err != nil {
		return err
	}

	return err
}

func checkStringExistsFile(filename, value string) (bool, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == value {
			return true, nil
		}
	}

	return false, scanner.Err()
}

func checkInstalled(targetDir string) bool {
	if _, err := os.Stat(filepath.Join(targetDir, unpackedOkay)); err == nil {
		return true
	}
	return false
}

func setInstalled(targetDir string) error {
	return os.WriteFile(filepath.Join(targetDir, unpackedOkay), nil, 0644)
}
