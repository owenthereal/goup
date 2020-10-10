package commands

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	goHost                = "golang.org"
	goDownloadBaseURL     = "https://dl.google.com/go"
	goSourceGitURL        = "https://github.com/golang/go"
	goSourceUpsteamGitURL = "https://go.googlesource.com/go"
)

var (
	installCmdGoHostFlag string

	installCmd = &cobra.Command{
		Use:   "install [version]",
		Short: `Install Go with a version`,
		Long: `Install Go by providing a version. If no version is provided, install
the latest Go. If the version is 'tip', an optional change list (CL)
number can be provided.`,
		Example: `
  goup install
  goup install 1.15.2
  goup install go1.15.2
  goup install tip # Compile Go tip
  goup install tip 1234 # 1234 is the CL number
`,
		PersistentPreRunE: preRunInstall,
		RunE:              runInstall,
	}
)

func init() {
	gh := os.Getenv("GOUP_GO_HOST")
	if gh == "" {
		gh = goHost
	}

	installCmd.PersistentFlags().StringVar(&installCmdGoHostFlag, "host", gh, "host that is used to download Go. The GOUP_GO_HOST environment variable overrides this flag.")
}

func preRunInstall(cmd *cobra.Command, args []string) error {
	http.DefaultTransport = &userAgentTransport{http.DefaultTransport}
	return nil
}

func runInstall(cmd *cobra.Command, args []string) error {
	var (
		ver string
		err error
	)

	if len(args) == 0 {
		ver, err = latestGoVersion()
		if err != nil {
			return err
		}
	} else {
		ver = args[0]
	}

	// Add go prefix, e.g., go1.15.2
	if !strings.HasPrefix(ver, "go") {
		ver = "go" + ver
	}

	if ver == "gotip" {
		var cl string
		if len(args) > 1 {
			cl = args[1]
		}

		err = installTip(cl)
	} else {
		err = install(ver)
	}
	if err != nil {
		return err
	}

	if err := symlink(ver); err != nil {
		return err
	}

	logger.Printf("Default Go is set to '%s'", ver)

	return nil
}

func latestGoVersion() (string, error) {
	h, err := url.Parse(installCmdGoHostFlag)
	if err != nil {
		return "", err
	}
	if h.Scheme == "" {
		h.Scheme = "https"
	}

	resp, err := http.Get(fmt.Sprintf("%s/VERSION?m=text", h.String()))
	if err != nil {
		return "", fmt.Errorf("Getting current Go version failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		b, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("Could not get current Go version: HTTP %d: %q", resp.StatusCode, b)
	}
	version, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(version)), nil
}

func symlink(ver string) error {
	current := currentGoRootDir()
	version := versionGoRootDir(ver)

	// ignore error, similar to rm -f
	os.Remove(current)

	return os.Symlink(version, current)
}

func install(version string) error {
	targetDir := versionGoRootDir(version)

	if _, err := os.Stat(filepath.Join(targetDir, unpackedOkay)); err == nil {
		logger.Printf("%s: already downloaded in %v", version, targetDir)
		return nil
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}
	goURL := versionArchiveURL(version)
	res, err := http.Head(goURL)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no binary release of %v for %v/%v at %v", version, getOS(), runtime.GOARCH, goURL)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %v checking size of %v", http.StatusText(res.StatusCode), goURL)
	}
	base := path.Base(goURL)
	archiveFile := filepath.Join(targetDir, base)
	if fi, err := os.Stat(archiveFile); err != nil || fi.Size() != res.ContentLength {
		if err != nil && !os.IsNotExist(err) {
			// Something weird. Don't try to download.
			return err
		}
		if err := copyFromURL(archiveFile, goURL); err != nil {
			return fmt.Errorf("error downloading %v: %v", goURL, err)
		}
		fi, err = os.Stat(archiveFile)
		if err != nil {
			return err
		}
		if fi.Size() != res.ContentLength {
			return fmt.Errorf("downloaded file %s size %v doesn't match server size %v", archiveFile, fi.Size(), res.ContentLength)
		}
	}
	wantSHA, err := slurpURLToString(goURL + ".sha256")
	if err != nil {
		return err
	}
	if err := verifySHA256(archiveFile, strings.TrimSpace(wantSHA)); err != nil {
		return fmt.Errorf("error verifying SHA256 of %v: %v", archiveFile, err)
	}
	logger.Printf("Unpacking %v ...", archiveFile)
	if err := unpackArchive(targetDir, archiveFile); err != nil {
		return fmt.Errorf("extracting archive %v: %v", archiveFile, err)
	}
	if err := ioutil.WriteFile(filepath.Join(targetDir, unpackedOkay), nil, 0644); err != nil {
		return err
	}
	logger.Printf("Success: %s downloaded in %v", version, targetDir)
	return nil
}

func installTip(clNumber string) error {
	root := versionGoRootDir("gotip")

	git := func(args ...string) error {
		cmd := exec.Command("git", args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = root
		return cmd.Run()
	}
	gitOutput := func(args ...string) ([]byte, error) {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		return cmd.Output()
	}

	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		if err := os.MkdirAll(root, 0755); err != nil {
			return fmt.Errorf("failed to create repository: %v", err)
		}
		if err := git("clone", "--depth=1", goSourceGitURL, root); err != nil {
			return fmt.Errorf("failed to clone git repository: %v", err)
		}
		if err := git("remote", "add", "upstream", goSourceUpsteamGitURL); err != nil {
			return fmt.Errorf("failed to add upstream git repository: %v", err)
		}
	}

	if clNumber != "" {
		fmt.Fprintf(os.Stderr, "This will download and execute code from golang.org/cl/%s, continue? [y/n] ", clNumber)
		var answer string
		if fmt.Scanln(&answer); answer != "y" {
			return fmt.Errorf("interrupted")
		}

		// CL is for googlesource, ls-remote against upstream
		// ls-remote outputs a number of lines like:
		// 2621ba2c60d05ec0b9ef37cd71e45047b004cead	refs/changes/37/227037/1
		// 51f2af2be0878e1541d2769bd9d977a7e99db9ab	refs/changes/37/227037/2
		// af1f3b008281c61c54a5d203ffb69334b7af007c	refs/changes/37/227037/3
		// 6a10ebae05ce4b01cb93b73c47bef67c0f5c5f2a	refs/changes/37/227037/meta
		refs, err := gitOutput("ls-remote", "upstream")
		if err != nil {
			return fmt.Errorf("failed to list remotes: %v", err)
		}
		r := regexp.MustCompile(`refs/changes/\d\d/` + clNumber + `/(\d+)`)
		match := r.FindAllStringSubmatch(string(refs), -1)
		if match == nil {
			return fmt.Errorf("CL %v not found", clNumber)
		}
		var ref string
		var patchSet int
		for _, m := range match {
			ps, _ := strconv.Atoi(m[1])
			if ps > patchSet {
				patchSet = ps
				ref = m[0]
			}
		}
		logger.Printf("Fetching CL %v, Patch Set %v...", clNumber, patchSet)
		if err := git("fetch", "upstream", ref); err != nil {
			return fmt.Errorf("failed to fetch %s: %v", ref, err)
		}
	} else {
		logger.Printf("Updating the go development tree...")
		if err := git("fetch", "origin", "master"); err != nil {
			return fmt.Errorf("failed to fetch git repository updates: %v", err)
		}
	}

	// Use checkout and a detached HEAD, because it will refuse to overwrite
	// local changes, and warn if commits are being left behind, but will not
	// mind if master is force-pushed upstream.
	if err := git("-c", "advice.detachedHead=false", "checkout", "FETCH_HEAD"); err != nil {
		return fmt.Errorf("failed to checkout git repository: %v", err)
	}
	// It shouldn't be the case, but in practice sometimes binary artifacts
	// generated by earlier Go versions interfere with the build.
	//
	// Ask the user what to do about them if they are not gitignored. They might
	// be artifacts that used to be ignored in previous versions, or precious
	// uncommitted source files.
	if err := git("clean", "-i", "-d"); err != nil {
		return fmt.Errorf("failed to cleanup git repository: %v", err)
	}
	// Wipe away probably boring ignored files without bothering the user.
	if err := git("clean", "-q", "-f", "-d", "-X"); err != nil {
		return fmt.Errorf("failed to cleanup git repository: %v", err)
	}

	cmd := exec.Command(filepath.Join(root, "src", makeScript()))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(root, "src")
	if runtime.GOOS == "windows" {
		// Workaround make.bat not autodetecting GOROOT_BOOTSTRAP. Issue 28641.
		goroot, err := exec.Command("go", "env", "GOROOT").Output()
		if err != nil {
			return fmt.Errorf("failed to detect an existing go installation for bootstrap: %v", err)
		}
		cmd.Env = append(os.Environ(), "GOROOT_BOOTSTRAP="+strings.TrimSpace(string(goroot)))
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build go: %v", err)
	}

	return nil
}

// unpackArchive unpacks the provided archive zip or tar.gz file to targetDir,
// removing the "go/" prefix from file entries.
func unpackArchive(targetDir, archiveFile string) error {
	switch {
	case strings.HasSuffix(archiveFile, ".zip"):
		return unpackZip(targetDir, archiveFile)
	case strings.HasSuffix(archiveFile, ".tar.gz"):
		return unpackTarGz(targetDir, archiveFile)
	default:
		return errors.New("unsupported archive file")
	}
}

// unpackTarGz is the tar.gz implementation of unpackArchive.
func unpackTarGz(targetDir, archiveFile string) error {
	r, err := os.Open(archiveFile)
	if err != nil {
		return err
	}
	defer r.Close()
	madeDir := map[string]bool{}
	zr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !validRelPath(f.Name) {
			return fmt.Errorf("tar file contained invalid name %q", f.Name)
		}
		rel := filepath.FromSlash(strings.TrimPrefix(f.Name, "go/"))
		abs := filepath.Join(targetDir, rel)

		fi := f.FileInfo()
		mode := fi.Mode()
		switch {
		case mode.IsRegular():
			// Make the directory. This is redundant because it should
			// already be made by a directory entry in the tar
			// beforehand. Thus, don't check for errors; the next
			// write will fail with the same error.
			dir := filepath.Dir(abs)
			if !madeDir[dir] {
				if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
					return err
				}
				madeDir[dir] = true
			}
			wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}
			n, err := io.Copy(wf, tr)
			if closeErr := wf.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
			if err != nil {
				return fmt.Errorf("error writing to %s: %v", abs, err)
			}
			if n != f.Size {
				return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, abs, f.Size)
			}
			if !f.ModTime.IsZero() {
				if err := os.Chtimes(abs, f.ModTime, f.ModTime); err != nil {
					// benign error. Gerrit doesn't even set the
					// modtime in these, and we don't end up relying
					// on it anywhere (the gomote push command relies
					// on digests only), so this is a little pointless
					// for now.
					logger.Printf("error changing modtime: %v", err)
				}
			}
		case mode.IsDir():
			if err := os.MkdirAll(abs, 0755); err != nil {
				return err
			}
			madeDir[abs] = true
		default:
			return fmt.Errorf("tar file entry %s contained unsupported file type %v", f.Name, mode)
		}
	}
	return nil
}

// unpackZip is the zip implementation of unpackArchive.
func unpackZip(targetDir, archiveFile string) error {
	zr, err := zip.OpenReader(archiveFile)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, f := range zr.File {
		name := strings.TrimPrefix(f.Name, "go/")

		outpath := filepath.Join(targetDir, name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outpath, 0755); err != nil {
				return err
			}
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		// File
		if err := os.MkdirAll(filepath.Dir(outpath), 0755); err != nil {
			return err
		}
		out, err := os.OpenFile(outpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		if err != nil {
			out.Close()
			return err
		}
		if err := out.Close(); err != nil {
			return err
		}
	}
	return nil
}

// verifySHA256 reports whether the named file has contents with
// SHA-256 of the given wantHex value.
func verifySHA256(file, wantHex string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return err
	}
	if fmt.Sprintf("%x", hash.Sum(nil)) != wantHex {
		return fmt.Errorf("%s corrupt? does not have expected SHA-256 of %v", file, wantHex)
	}
	return nil
}

// slurpURLToString downloads the given URL and returns it as a string.
func slurpURLToString(url_ string) (string, error) {
	res, err := http.Get(url_)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %v", url_, res.Status)
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %v", url_, err)
	}
	return string(slurp), nil
}

// copyFromURL downloads srcURL to dstFile.
func copyFromURL(dstFile, srcURL string) (err error) {
	f, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			f.Close()
			os.Remove(dstFile)
		}
	}()
	c := &http.Client{
		Transport: &userAgentTransport{&http.Transport{
			// It's already compressed. Prefer accurate ContentLength.
			// (Not that GCS would try to compress it, though)
			DisableCompression: true,
			DisableKeepAlives:  true,
			Proxy:              http.ProxyFromEnvironment,
		}},
	}
	res, err := c.Get(srcURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	pw := &progressWriter{w: f, total: res.ContentLength}
	n, err := io.Copy(pw, res.Body)
	if err != nil {
		return err
	}
	if res.ContentLength != -1 && res.ContentLength != n {
		return fmt.Errorf("copied %v bytes; expected %v", n, res.ContentLength)
	}
	pw.update() // 100%
	return f.Close()
}

func makeScript() string {
	switch runtime.GOOS {
	case "plan9":
		return "make.rc"
	case "windows":
		return "make.bat"
	default:
		return "make.bash"
	}
}

type progressWriter struct {
	w     io.Writer
	n     int64
	total int64
	last  time.Time
}

func (p *progressWriter) update() {
	end := " ..."
	if p.n == p.total {
		end = ""
	}
	fmt.Fprintf(os.Stderr, "Downloaded %5.1f%% (%*d / %d bytes)%s\n",
		(100.0*float64(p.n))/float64(p.total),
		ndigits(p.total), p.n, p.total, end)
}

func ndigits(i int64) int {
	var n int
	for ; i != 0; i /= 10 {
		n++
	}
	return n
}

func (p *progressWriter) Write(buf []byte) (n int, err error) {
	n, err = p.w.Write(buf)
	p.n += int64(n)
	if now := time.Now(); now.Unix() != p.last.Unix() {
		p.update()
		p.last = now
	}
	return
}

// getOS returns runtime.GOOS. It exists as a function just for lazy
// testing of the Windows zip path when running on Linux/Darwin.
func getOS() string {
	return runtime.GOOS
}

// versionArchiveURL returns the zip or tar.gz URL of the given Go version.
func versionArchiveURL(version string) string {
	goos := getOS()

	ext := "tar.gz"
	if goos == "windows" {
		ext = "zip"
	}

	arch := runtime.GOARCH
	if goos == "linux" && runtime.GOARCH == "arm" {
		arch = "armv6l"
	}

	return fmt.Sprintf("%s/%s.%s-%s.%s", goDownloadBaseURL, version, goos, arch, ext)
}

// unpackedOkay is a sentinel zero-byte file to indicate that the Go
// version was downloaded and unpacked successfully.
const unpackedOkay = ".unpacked-success"

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

type userAgentTransport struct {
	rt http.RoundTripper
}

func (uat userAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	version := runtime.Version()
	if strings.Contains(version, "devel") {
		// Strip the SHA hash and date. We don't want spaces or other tokens (see RFC2616 14.43)
		version = "devel"
	}
	r.Header.Set("User-Agent", "goup/"+version)
	return uat.rt.RoundTrip(r)
}
