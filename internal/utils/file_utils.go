package utils

import "os"

func CreateDataFolder(baseDir string) {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		os.MkdirAll(baseDir, 0777)
	}
}
