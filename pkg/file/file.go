package file

import (
	"io/ioutil"
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

func WalkDir(dir string) ([]string, error) {
	var fileList []string
	if files, err := ioutil.ReadDir(dir); err != nil {
		return nil, err
	} else {
		for _, file := range files {
			path := filepath.Join(dir, file.Name())
			if file.IsDir() {
				if subDirFiles, err := WalkDir(path); err != nil {
					return nil, err
				} else {
					fileList = append(fileList, subDirFiles...)
				}
			} else {
				fileList = append(fileList, path)
			}
		}
	}
	return fileList, nil
}
