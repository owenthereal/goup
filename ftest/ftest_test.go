package ftest

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/owenthereal/goup/internal/commands"
)

var (
	goupBinDir string
	goupBin    string

	flagMainPath    string
	flagInstallPath string
)

func TestMain(m *testing.M) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&flagMainPath, "main-path", filepath.Join(pwd, "..", "cmd", "goup"), "main.go path")
	flag.StringVar(&flagInstallPath, "install-path", filepath.Join(pwd, "..", "install.sh"), "install.sh path")
	flag.Parse()

	goupBinDir, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(goupBinDir)

	goupBin = filepath.Join(goupBinDir, runtime.GOOS+"-"+runtime.GOARCH)
	cmd := exec.Command("go", "build", "-o", goupBin, flagMainPath)

	log.Printf("%s", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s: %s", out, err)
	}

	os.Exit(m.Run())
}

func TestInstaller(t *testing.T) {
	cmd := exec.Command("sh", flagInstallPath, "--skip-prompt")
	cmd.Env = append(os.Environ(), "GOUP_UPDATE_ROOT=file://"+goupBinDir)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %s", out, err)
	}

	fmt.Println(string(out))

	// check file exists
	filesShouldExist := []string{
		commands.GoupDir(),
		commands.GoupEnvFile(),
		commands.GoupBinDir(),
		commands.GoupCurrentDir(),
	}
	for _, f := range filesShouldExist {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Error(err)
		}
	}

	// check profiles
	for _, f := range commands.ProfileFiles {
		ok, err := fileContains(f, commands.ProfileFileSourceContent)
		if err != nil {
			t.Error(err)
		}

		if !ok {
			t.Errorf("%s does not source goup", f)
		}
	}
}

func fileContains(f string, s string) (bool, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return false, err
	}

	return bytes.Contains(b, []byte(s)), nil
}
