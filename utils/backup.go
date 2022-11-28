package utils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

func BackupBuildFile(path string) string {
	tmpFile, err := ioutil.TempFile(".", "BUILD-*.bak")

	if err != nil {
		log.Fatalln(err)
	}

	defer tmpFile.Close()

	srcFile, err := os.Open(path)

	if err != nil {
		log.Fatalln(err)
	}

	defer srcFile.Close()

	io.Copy(srcFile, tmpFile)
	tmpFile.Sync()

	return tmpFile.Name()
}

func RestoreBackup(backup string, orig string) {
	backupFile, _ := os.Open(backup)
	defer backupFile.Close()
	defer os.Remove(backupFile.Name())

	origFile, _ := os.Open(orig)
	defer origFile.Close()

	io.Copy(backupFile, origFile)
}
