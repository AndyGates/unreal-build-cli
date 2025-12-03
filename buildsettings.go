package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type BuildSettings struct {
	ClientPlatforms    []string
	ServerPlatforms    []string
	Configurations     []string
	Steps              []string
	CookType           CookType
	AdditionalCookArgs []string
}

func GetPresetPath(presetName string) (string, error) {

	saveDir, err := GetSaveDirectory()
	if err != nil {
		log.Printf("Couldn't get save dir, preset will not be saved: %s", err)
		return "", err
	}

	presetFile := fmt.Sprintf("%s.json", presetName)
	return filepath.Join(saveDir, presetFile), nil
}

func GetPresetList() ([]string, error) {

	var matchingFiles []string

	saveDir, err := GetSaveDirectory()
	if err != nil {
		log.Printf("Couldn't get save dir, cannot list presets: %s", err)
		return matchingFiles, err
	}

	pattern := "*.json"
	files, err := os.ReadDir(saveDir)
	if err != nil {
		fmt.Println("Could not enumerate save dir:", err)
		return matchingFiles, err
	}
	for _, file := range files {
		match, err := filepath.Match(pattern, file.Name())
		if err != nil {
			fmt.Println("Error matching file in save dir:", err)
			continue
		}

		if match {
			filename := file.Name()
			matchingFiles = append(matchingFiles, strings.TrimSuffix(filename, filepath.Ext(filename)))
		}
	}

	return matchingFiles, nil
}

func SaveBuildSettings(settings BuildSettings) {

	settingsJson, _ := json.MarshalIndent(settings, "", "    ")

	savePath, err := GetPresetPath("last")
	if err != nil {
		log.Printf("Not saving preset: %s", err)
		return
	}

	saveDir := filepath.Dir(savePath)

	if !CheckFileExists(saveDir) {
		err := os.MkdirAll(saveDir, 0700) // Create your file

		if err != nil {
			log.Printf("Save path did not exist and we could not create it: %s", err)
			return
		}
	}

	err = os.WriteFile(savePath, settingsJson, fs.ModePerm)

	if err != nil {
		log.Printf("Failed to write preset: %s", err)
	}
}

func LoadBuildSettings(name string) (BuildSettings, error) {

	loadedBuildSettings := BuildSettings{}

	lastPresetPath, err := GetPresetPath(name)
	if err != nil {
		return loadedBuildSettings, os.ErrNotExist
	}

	if !CheckFileExists(lastPresetPath) {
		return loadedBuildSettings, os.ErrNotExist
	}

	file, _ := os.ReadFile(lastPresetPath)
	err = json.Unmarshal(file, &loadedBuildSettings)

	return loadedBuildSettings, err
}
