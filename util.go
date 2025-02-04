package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqroot/prompt"
	"golang.org/x/sys/windows/registry"
)

const ProjectSaveDir = "unreal-build-cli"
const EngineAssociationRegistryKey = `SOFTWARE\Epic Games\Unreal Engine\Builds`
const RelativeRunUATPath = `/Engine/Build/BatchFiles/RunUAT.bat`

func queryRegistryKeyValue(path string, name string) (string, error) {
	// Open the registry key
	key, err := registry.OpenKey(registry.CURRENT_USER, path, registry.READ)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err := key.GetStringValue(name)
	if err != nil {
		return "", err
	}

	return value, nil
}

func getUprojectEngineAssociation(uprojectPath string) (string, error) {
	// Open the uproject file
	uprojectFile, err := os.Open(uprojectPath)
	if err != nil {
		log.Fatal(err)
	}
	defer uprojectFile.Close()

	byteValue, err := io.ReadAll(uprojectFile)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the uproject into a map
	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	// try get the engine association
	if value, found := result["EngineAssociation"]; found {
		return value.(string), nil
	} else {
		return "", errors.New("EngineAssociation not found")
	}
}

func getEngineDir(uprojectPath string) (string, error) {
	assoc, err := getUprojectEngineAssociation(uprojectPath)
	if err != nil {
		return "", err
	}

	// if this is a relative path, return it directly
	if strings.HasPrefix(assoc, ".") {
		return assoc, nil
	}

	// if this is a guid association, look it up in the registry
	if strings.HasPrefix(assoc, "{") {
		path, err := queryRegistryKeyValue(EngineAssociationRegistryKey, assoc)
		if err != nil {
			return "", err
		}
		return path, nil
	}

	//TODO: implement numbered vanilla installs (e.g. "5.3")
	return "", errors.New("unsupported engine association")
}

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

func GetRunUATPath(uprojectPath string) (string, error) {

	engineDir, err := getEngineDir(uprojectPath)
	if err != nil {
		return "", err
	}

	return filepath.Join(engineDir, RelativeRunUATPath), nil
}

func GetSaveDirectory() (string, error) {

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, ProjectSaveDir), err
}
