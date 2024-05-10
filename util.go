package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cqroot/prompt"
)

const ProjectSaveDir = "unreal-build-cli"

func CheckErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}

func Contains(array []string, match string) bool {
	for _, s := range array {
		if s == match {
			return true
		}
	}
	return false
}

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func FindUproject() string {

	dir := "."

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".uproject" {
			return filepath.Join(path, f.Name())
		}
	}

	return ""
}

func GetSaveDirectory() (string, error) {

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, ProjectSaveDir), err
}
