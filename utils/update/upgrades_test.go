package update

import (
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/requests"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUpgradeCli(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	tempDir := os.TempDir()
	version, err := requests.GetLatestVersion("stable")
	if err != nil {
		logrus.WithError(err).Debugln("get latest version")
		return
	}

	if version == info.Version {
		return
	}

	err = DownloadBinaryCompressedFile(
		GetBinaryCompressedFileName(version),
		tempDir,
	)
	if err != nil {
		return
	}

	f, err := os.Open(filepath.Join(tempDir, GetBinaryCompressedFileName(version)))
	if err != nil {
		logrus.WithError(err).Debugln("get compressed file")
		return
	}

	err = Untar(tempDir, f)
	if err != nil {
		logrus.WithError(err).Debugln("untar compressed file")
		return
	}

	// add permission
	err = os.Chmod(filepath.Join(tempDir, "lets"), 0755)
	if err != nil {
		logrus.WithError(err).Debugln("get compressed file")
		return
	}

	// replace binary
	path, err := exec.LookPath("lets")
	if err != nil {
		logrus.WithError(err).Debugln("get compressed file")
		return
	}

	err = os.Rename(filepath.Join(tempDir, "lets"), path)
	if err != nil {
		logrus.WithError(err).Debugln("get compressed file")
		return
	}
}
