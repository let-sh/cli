package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetFilesSize(paths []string) (int64, error) {
	var size int64
	for _, path := range paths {
		s, err := fileSize(path)
		if err != nil {
			return 0, err
		}
		size += s
	}
	return size, nil
}

func fileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

func Move(src, dst string) error {
	var err error
	// This returns an *os.FileInfo type
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		err = MoveDirectory(src, dst)
	} else {
		err = MoveFile(src, dst)
	}
	return err
}

func MoveDirectory(src, dst string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, file := range files {
		srcfp := filepath.Join(src, file.Name())
		dstfp := filepath.Join(dst, file.Name())
		if file.IsDir() {
			MoveDirectory(srcfp, dstfp)
		} else {
			MoveFile(srcfp, dstfp)
		}
	}
	return nil
}

func MoveFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
