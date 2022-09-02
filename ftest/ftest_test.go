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

func TestGoup(t *testing.T) {
	t.Run("installer", func(t *testing.T) {
		cmd := exec.Command("sh", flagInstallPath, "--skip-prompt")
		cmd.Env = append(os.Environ(), "GOUP_UPDATE_ROOT=file://"+goupBinDir)
		execCmd(t, cmd)

		// check file exists
		filesShouldExist := []string{
			commands.GoupDir(),
			commands.GoupEnvFile(),
			commands.GoupBinDir(),
			commands.GoupCurrentDir(),
			commands.GoupCurrentBinDir(),
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
	})

	goupBin := filepath.Join(commands.GoupBinDir(), "goup")

	t.Run("goup install 1.15.2", func(t *testing.T) {
		cmd := exec.Command(goupBin, "install", "1.15.2")
		execCmd(t, cmd)
	})

	t.Run("goup show 1.15.2", func(t *testing.T) {
		cmd := exec.Command(goupBin, "show")
		out := execCmd(t, cmd)

		if want, got := []byte("1.15.2"), out; !bytes.Contains(got, want) {
			t.Fatalf("goup show failed: want=%s got=%s", want, out)
		}
	})

	t.Run("goup install 1.15.3", func(t *testing.T) {
		cmd := exec.Command(goupBin, "install", "1.15.3")
		execCmd(t, cmd)
	})

	t.Run("goup show 1.15.3", func(t *testing.T) {
		cmd := exec.Command(goupBin, "show")
		out := execCmd(t, cmd)

		if want, got := []byte("1.15.3"), out; !bytes.Contains(got, want) {
			t.Fatalf("goup show failed: want=%s got=%s", want, out)
		}
	})

	t.Run("goup default 1.15.2", func(t *testing.T) {
		cmd := exec.Command(goupBin, "default", "1.15.2")
		execCmd(t, cmd)
	})

	t.Run("goup show 1.15.2", func(t *testing.T) {
		cmd := exec.Command(goupBin, "show")
		out := execCmd(t, cmd)

		if want, got := []byte("1.15.2"), out; !bytes.Contains(got, want) {
			t.Fatalf("goup show failed: want=%s got=%s", want, out)
		}
	})

	t.Run("goup search", func(t *testing.T) {
		cmd := exec.Command(goupBin, "search")
		out := execCmd(t, cmd)

		if want, got := []byte("1.15.2"), out; !bytes.Contains(got, want) {
			t.Fatalf("goup search failed: want=%s got=%s", want, out)
		}
	})

	t.Run("goup remove 1.15.2", func(t *testing.T) {
		cmd := exec.Command(goupBin, "remove", "1.15.2")
		execCmd(t, cmd)
	})

	t.Run("goup show does not have 1.15.2", func(t *testing.T) {
		cmd := exec.Command(goupBin, "show")
		out := execCmd(t, cmd)

		if want, got := []byte("1.15.2"), out; bytes.Contains(got, want) {
			t.Fatalf("goup show again failed: want=%s got=%s", want, out)
		}
	})

	t.Run("goup upgrade", func(t *testing.T) {
		cmd := exec.Command(goupBin, "upgrade", "0.1.4")
		cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s:$PATH", filepath.Dir(goupBin)))
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("goup upgrade failed: %s: %s", out, err)
		}
		t.Logf("%s", out)
	})

	t.Run("goup version", func(t *testing.T) {
		cmd := exec.Command(goupBin, "version")
		out := execCmd(t, cmd)

		if want, got := []byte("0.1.4"), out; !bytes.Contains(got, want) {
			t.Fatalf("goup version failed: want=%s got=%s", want, out)
		}
	})

}

func execCmd(t *testing.T, cmd *exec.Cmd) []byte {
	t.Helper()

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %s", out, err)
	}
	t.Logf("%s", out)

	return out
}

func fileContains(f string, s string) (bool, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return false, err
	}

	return bytes.Contains(b, []byte(s)), nil
}
