package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/let-sh/cli/info"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func UpgradeCli(force bool, channel string) {
	binaryName := "lets"
	if runtime.GOOS == "windows" {
		binaryName = "lets.exe"
	}

	tempDir := os.TempDir()
	logrus.Debugf("tempDir: %s", tempDir)

	// if not specified channel, using current channel
	if channel == "" {
		channel = GetCurrentReleaseChannel()
	}

	if channel != "dev" || force {
		version, err := requests.GetLatestVersion(channel)
		if err != nil {
			log.Warning("upgrade failed: " + err.Error())
			logrus.WithError(err).Debugln("get latest version")
			return
		}

		if version == info.Version && !force {
			log.Success("currently is the latest version: " + info.Version)
			return
		}

		err = DownloadBinaryCompressedFile(
			GetBinaryCompressedFileName(version),
			tempDir,
		)
		if err != nil {
			log.Warning("upgrade failed: " + err.Error())
			logrus.WithError(err).Debugln("download binary")
			return
		}

		logrus.Debugf("open compressed file: %s", filepath.Join(tempDir, GetBinaryCompressedFileName(version)))
		f, err := os.Open(filepath.Join(tempDir, GetBinaryCompressedFileName(version)))
		if err != nil {
			log.Warning("upgrade failed: " + err.Error())
			logrus.WithError(err).Debugln("get compressed file")
			return
		}
		logrus.Debugf("compressed file: %s", f.Name())

		unzipedDir := filepath.Join(tempDir, "let.sh")
		os.Mkdir(unzipedDir, 0755)
		if runtime.GOOS == "windows" {
			err = Unzip(filepath.Join(tempDir, GetBinaryCompressedFileName(version)), unzipedDir)
		} else {
			err = Untar(unzipedDir, f)
		}

		if err != nil {
			log.Warning("upgrade failed: " + err.Error())
			logrus.WithError(err).Debugln("get compressed file")
			return
		}

		// add permission
		err = os.Chmod(filepath.Join(unzipedDir, binaryName), 0755)
		if err != nil {
			logrus.WithError(err).Debugln("get compressed file")
			return
		}
		logrus.Debugf("add permission: %s", filepath.Join(unzipedDir, binaryName))

		// replace binary
		path, err := exec.LookPath(binaryName)
		if err != nil {
			log.Warning("upgrade failed: " + err.Error())
			logrus.WithError(err).Debugln("get compressed file")
			return
		}

		if runtime.GOOS == "windows" {
			// handle process error
			// go further: https://stackoverflow.com/questions/9162969/how-can-a-c-binary-replace-itself
			err = os.Rename(path, path+".old")
			err = os.Rename(filepath.Join(unzipedDir, binaryName), path)
			os.RemoveAll(path + ".old")
		} else {
			err = os.Rename(filepath.Join(unzipedDir, binaryName), path)
		}

		if err != nil {
			log.Warning("upgrade failed: " + err.Error())
			logrus.WithError(err).Debugln("get compressed file")
			return
		}
		logrus.Debugf("mv binary: %s -> %s", filepath.Join(tempDir, binaryName), path)

		// clean up temporary files
		os.RemoveAll(unzipedDir)
		os.RemoveAll(filepath.Join(tempDir, GetBinaryCompressedFileName(version)))

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
	resp, err := http.Get("https://install.let-sh.com/" + filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"status_code": resp.Status,
			"url":         "https://install.let-sh.com/" + filename,
		}).WithError(err).Debugln("download binary compressed file error")
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	bodySize, err := strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	if err != nil {
		return fmt.Errorf("error requests")
	}

	p := mpb.New(
		mpb.WithWidth(64),
		mpb.WithRefreshRate(200*time.Millisecond),
	)

	bar := p.AddBar(bodySize,
		//mpb.BarStyle("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKiloByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKB, "% .2f", 1024),
		),
		mpb.BarRemoveOnComplete(),
	)

	// create proxy reader
	proxyReader := bar.ProxyReader(resp.Body)
	defer proxyReader.Close()

	// copy from proxyReader, ignoring errors
	io.Copy(out, proxyReader)

	// Writer the body to file
	//_, err = ui.Copy(out, resp.Body)
	//if err != nil {
	//	return err
	//}

	return nil
}

func GetBinaryCompressedFileName(version string) string {
	if runtime.GOOS == "windows" {
		return "cli_" + version + "_" + runtime.GOOS + "_" + runtime.GOARCH + ".zip"
	}
	return "cli_" + version + "_" + runtime.GOOS + "_" + runtime.GOARCH + ".tar.gz"
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
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
