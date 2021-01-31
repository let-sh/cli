package update

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func UpgradeCli() {
	tempDir := os.TempDir()
	logrus.Debugf("tempDir: %s", tempDir)
	if GetCurrentReleaseChannel() != "dev" {
		version, err := requests.GetLatestVersion(GetCurrentReleaseChannel())
		if err != nil {
			logrus.WithError(err).Debugln("get latest version")
			return
		}

		if version == info.Version {
			log.Success("currently is the latest version: " + info.Version)
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
		logrus.Debugf("compressed file: %s", f.Name())

		err = Untar(filepath.Join(tempDir, "lets"), f)
		if err != nil {
			logrus.WithError(err).Debugln("get compressed file")
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
		log.Success(fmt.Sprintf("Successfully installed let.sh %s!", version))

	} else {
		log.Warning("currently at development channel, no need to upgrade")
	}
}

func DownloadBinaryCompressedFile(filename, tempDir string) error {
	localFile := filepath.Join(tempDir, filename)

	// Create the file
	out, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get("https://install.let.sh.cn/" + filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"status_code": resp.Status,
			"url":         "https://install.let.sh.cn/" + filename,
		}).WithError(err).Debugln("download binary compressed file error")
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetBinaryCompressedFileName(version string) string {
	return "cli_" + version + "_" + runtime.GOOS + "_" + runtime.GOARCH + ".tar.gz"
}

func Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
