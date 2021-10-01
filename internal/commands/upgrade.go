package commands

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/kierdavis/ansi"
	"github.com/spf13/cobra"
	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
	"golang.org/x/oauth2"
)

func upgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade [version]",
		Short: "Upgrade goup",
		Long:  "Upgrade goup by providing a version. If no version is provided, upgrade to the latest goup.",
		RunE:  runUpgrade,
	}
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	ansi.HideCursor()
	defer ansi.ShowCursor()

	m := &manager{
		Manager: &update.Manager{
			Command: "goup",
			Store: &githubStore{
				Owner:       "owenthereal",
				Repo:        "goup",
				Version:     Version,
				AccessToken: os.Getenv("GITHUB_TOKEN"),
			},
		},
	}

	var r release
	if len(args) > 0 {
		rr, err := m.GetRelease(trimVPrefix(args[0]))
		if err != nil {
			return fmt.Errorf("error fetching release: %s", err)
		}

		r = release{rr}
	} else {
		// fetch the new releases
		releases, err := m.LatestReleases()
		if err != nil {
			log.Fatalf("error fetching releases: %s", err)
		}

		// no updates
		if len(releases) == 0 {
			logger.Println("No upgrades")
			return nil
		}

		// latest release
		r = release{releases[0]}

	}

	// find the tarball for this system
	a := r.FindTarballWithVersion(runtime.GOOS, runtime.GOARCH)
	if a == nil {
		return fmt.Errorf("no upgrade for your system")
	}

	bin, err := a.DownloadProxy(progress.Reader)
	if err != nil {
		return fmt.Errorf("error downloading: %s", err)
	}

	logger.Debugf("Downloaded release to %s", bin)

	// install it
	if err := m.InstallBin(bin); err != nil {
		return fmt.Errorf("error installing: %s", err)
	}

	logger.Printf("Upgraded to %s", trimVPrefix(r.Version))

	return nil
}

type manager struct {
	*update.Manager
}

func (m *manager) InstallBin(bin string) error {
	oldbin, err := exec.LookPath(m.Command)
	if err != nil {
		return fmt.Errorf("error looking up path of %q: %w", m.Command, err)
	}

	dir := filepath.Dir(oldbin)

	if err := os.Chmod(bin, 0755); err != nil {
		return fmt.Errorf("error in chmod: %w", err)
	}

	dst := filepath.Join(dir, m.Command)
	tmp := dst + ".tmp"

	logger.Debugf("Copy %q to %q", bin, tmp)
	if err := copyFile(tmp, bin); err != nil {
		return fmt.Errorf("error in copying: %w", err)
	}

	if runtime.GOOS == "windows" {
		old := dst + ".old"
		logger.Debugf("Windows workaround renaming %q to %q", dst, old)
		if err := os.Rename(dst, old); err != nil {
			return fmt.Errorf("error in windows renmaing: %w", err)
		}
	}

	logger.Debugf("Renaming %q to %q", tmp, dst)
	if err := os.Rename(tmp, dst); err != nil {
		return fmt.Errorf("error in renaming: %w", err)
	}

	return nil
}

type release struct {
	*update.Release
}

func (r *release) FindTarballWithVersion(os, arch string) *update.Asset {
	s := fmt.Sprintf("%s-%s", os, arch)
	for _, a := range r.Assets {
		if s == a.Name {
			return a
		}
	}

	return nil
}

func trimVPrefix(s string) string {
	return strings.TrimPrefix(s, "v")
}

// copyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func copyFile(dst, src string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}

	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

type githubStore struct {
	Owner       string
	Repo        string
	Version     string
	AccessToken string
}

func (s *githubStore) client(ctx context.Context) *github.Client {
	var tc *http.Client

	if token := s.AccessToken; token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc = oauth2.NewClient(ctx, ts)
	}

	return github.NewClient(tc)

}

// GetRelease returns the specified release or ErrNotFound.
func (s *githubStore) GetRelease(version string) (*update.Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gh := s.client(ctx)

	r, res, err := gh.Repositories.GetReleaseByTag(ctx, s.Owner, s.Repo, "v"+version)

	if res.StatusCode == 404 {
		return nil, update.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return githubRelease(r), nil
}

// LatestReleases returns releases newer than Version, or nil.
func (s *githubStore) LatestReleases() (latest []*update.Release, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gh := s.client(ctx)

	releases, _, err := gh.Repositories.ListReleases(ctx, s.Owner, s.Repo, nil)
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		tag := r.GetTagName()

		if tag == s.Version || "v"+s.Version == tag {
			break
		}

		latest = append(latest, githubRelease(r))
	}

	return
}

// githubRelease returns a Release.
func githubRelease(r *github.RepositoryRelease) *update.Release {
	out := &update.Release{
		Version:     r.GetTagName(),
		Notes:       r.GetBody(),
		PublishedAt: r.GetPublishedAt().Time,
		URL:         r.GetURL(),
	}

	for _, a := range r.Assets {
		out.Assets = append(out.Assets, &update.Asset{
			Name:      a.GetName(),
			Size:      a.GetSize(),
			URL:       a.GetBrowserDownloadURL(),
			Downloads: a.GetDownloadCount(),
		})
	}

	return out
}
