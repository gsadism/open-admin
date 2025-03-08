package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Exists  whether the given path exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Abs(dir string) (string, error) {
	if !filepath.IsAbs(dir) {
		if d, err := filepath.Abs(dir); err != nil {
			return "", err
		} else {
			if Exists(d) {
				return d, nil
			} else {
				return "", errors.New(fmt.Sprintf("%s is not exists", dir))
			}
		}
	} else {
		if Exists(dir) {
			return dir, nil
		} else {
			return "", errors.New(fmt.Sprintf("%s is not exists", dir))
		}
	}
}

func Folder(dir string) (string, error) {
	if !filepath.IsAbs(dir) {
		if d, err := filepath.Abs(dir); err != nil {
			return "", err
		} else {
			dir = d
		}
	}
	if !Exists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", err
		}
	}
	return dir, nil
}
