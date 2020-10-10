package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

const (
	envFileContent = `export PATH="$HOME/.go/bin":"$HOME/.go/current/bin":$PATH`
	sourceContent  = `source "$HOME/.go/env"`

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
	initInstallFlag bool
	initSkipPrompt  bool

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the goup environment file",
		RunE:  runInit,
	}
)

func init() {
	initCmd.PersistentFlags().BoolVar(&initInstallFlag, "install", false, "Install the latest Go")
	initCmd.PersistentFlags().BoolVar(&initSkipPrompt, "skip-prompt", false, "Skip confirmation prompt")
}

func runInit(cmd *cobra.Command, args []string) error {
	profiles := []string{
		filepath.Join(homedir, ".profile"),
		filepath.Join(homedir, ".zprofile"),
		filepath.Join(homedir, ".bash_profile"),
	}

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
		GoupDir:         goupDir(),
		GoupBinDir:      goupBinDir(),
		CurrentGoBinDir: currentGoBinDir(),
		ProfileFiles:    profiles,
	}
	if err := tmpl.Execute(os.Stdout, params); err != nil {
		return err
	}

	answer, err := prompt("\nWould you like to proceed with the installation?", "yes")
	if err != nil {
		return err
	}

	if answer != "yes" {
		return nil
	}

	// add a line break
	fmt.Println("")

	ef := envFile()
	if err := os.MkdirAll(filepath.Dir(ef), 0755); err != nil {
		return err
	}

	// ignore error, similar to rm -f
	os.Remove(ef)

	if err := ioutil.WriteFile(ef, []byte(envFileContent), 0664); err != nil {
		return err
	}

	if err := appendSourceToProfiles(profiles); err != nil {
		return err
	}

	if initInstallFlag {
		if err := runInstall(cmd, args); err != nil {
			return err
		}
	}

	return nil
}

func prompt(query, defaultAnswer string) (string, error) {
	if initSkipPrompt {
		return defaultAnswer, nil
	}

	fmt.Printf("%s [%s]: ", query, defaultAnswer)

	s := bufio.NewScanner(os.Stdin)
	if !s.Scan() {
		return "", s.Err()
	}

	answer := s.Text()
	if answer == "" {
		answer = defaultAnswer
	}

	return answer, nil
}

func appendSourceToProfiles(profiles []string) error {
	for _, profile := range profiles {
		if err := appendToFile(profile, sourceContent); err != nil {
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
