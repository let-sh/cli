package utils

import (
	"os"
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
