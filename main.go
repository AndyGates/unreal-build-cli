package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/multichoose"
)

func PromptForBuildSettings(config Config) BuildSettings {

	clientPlatforms, err := prompt.New().Ask("Select client platforms to build:").
		MultiChoose(config.ClientOptions.Options, multichoose.WithDefaultIndexes(0, config.ClientOptions.Defaults))
	CheckErr(err)

	serverPlatforms, err := prompt.New().Ask("Select server platforms to build:").
		MultiChoose(config.ServerOptions.Options, multichoose.WithDefaultIndexes(0, config.ServerOptions.Defaults))
	CheckErr(err)

	configurations, err := prompt.New().Ask("Select configurations to build:").
		MultiChoose(config.ConfigurationOptions.Options, multichoose.WithDefaultIndexes(0, config.ConfigurationOptions.Defaults))
	CheckErr(err)

	steps, err := prompt.New().Ask("Select build steps to run:").
		MultiChoose(config.StepOptions.Options, multichoose.WithDefaultIndexes(0, config.StepOptions.Defaults))
	CheckErr(err)

	shouldIterate := false

	if Contains(steps, config.StepOptions.Options[1]) {

		const yes = "yes"
		const no = "no"

		iterate, err := prompt.New().Ask("Run an iterative cook?").
			Choose([]string{yes, no},
				choose.WithTheme(choose.ThemeLine),
				choose.WithKeyMap(choose.HorizontalKeyMap))
		CheckErr(err)

		shouldIterate = iterate == yes
	}

	return BuildSettings{
		ClientPlatforms: clientPlatforms,
		ServerPlatforms: serverPlatforms,
		Configurations:  configurations,
		Steps:           steps,
		ShouldIterate:   shouldIterate,
	}
}

func PromptForPresets() (string, error) {

	presets, err := GetPresetList()
	if err != nil {
		return "", err
	}

	presetChoice, err := prompt.New().Ask("Select a preset to build:").
		Choose(presets)
	CheckErr(err)

	return presetChoice, nil
}

func buildArgumentList(uprojectPath string, config Config, buildSettings BuildSettings) []string {

	// RunUAT.bat BuildCookRun -project="..."
	args := []string{"BuildCookRun"}

	// -project
	args = append(args, fmt.Sprintf("-project=\"%s\"", uprojectPath))

	if len(buildSettings.ClientPlatforms) > 0 {
		// -targetplatform
		args = append(args, fmt.Sprintf("-targetplatform=%s", strings.Join(buildSettings.ClientPlatforms, "+")))
		// -configuration
		args = append(args, fmt.Sprintf("-clientconfig=%s", strings.Join(buildSettings.Configurations, "+")))
	}

	// -serverplatform
	if len(buildSettings.ServerPlatforms) > 0 {
		args = append(args, "-server")
		// -serverplatform
		args = append(args, fmt.Sprintf("-serverplatform=%s", strings.Join(buildSettings.ServerPlatforms, "+")))
		// -configuration
		args = append(args, fmt.Sprintf("-serverconfig=%s", strings.Join(buildSettings.Configurations, "+")))
	}

	// -cook -skipcook etc.
	for _, stepOption := range config.StepOptions.Options {
		stepArg := strings.ToLower(stepOption)
		if !Contains(buildSettings.Steps, stepOption) {
			stepArg = fmt.Sprintf("skip%s", stepArg)
		}
		args = append(args, fmt.Sprintf("-%s", stepArg))
	}

	// -iterate
	if buildSettings.ShouldIterate {
		args = append(args, "-iterate")
	}

	return args
}

func main() {
	config := GetConfig()
	uprojectPath := FindUproject()

	shouldReplayPtr := flag.Bool("r", false, "replay the last build settings")
	usePresetPtr := flag.Bool("p", false, "present a list of presets")

	flag.Parse()

	if uprojectPath == "" {
		log.Fatal("No .uproject found, make sure you are in the project directory")
	} else {
		log.Printf("Found project file: %s", uprojectPath)
	}

	var buildSettings *BuildSettings = nil

	if *shouldReplayPtr {
		loadedBuildSettings, err := LoadBuildSettings("last")
		if err != nil {
			buildSettings = &loadedBuildSettings
		}
	}

	if *usePresetPtr {
		preset, err := PromptForPresets()
		if err != nil {
			log.Printf("Failed to prompt for presets")
		}

		if preset != "" {
			loadedBuildSettings, err := LoadBuildSettings(preset)
			if err != nil {
				buildSettings = &loadedBuildSettings
			}
		}
	}

	if buildSettings == nil {
		promptedBuildSettings := PromptForBuildSettings(config)
		buildSettings = &promptedBuildSettings
	}

	args := buildArgumentList(uprojectPath, config, *buildSettings)

	log.Printf("Invoking RunUAT with arguments: %s", args)

	SaveBuildSettings(*buildSettings)

	cmd := exec.Command(config.RunUATPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println("could not run command: ", err)
	}
}

// ..\UnrealEngine\Engine\Build\BatchFiles\RunUAT.bat BuildCookRun -project="..\Yakisoba\Yakisoba.uproject" -platform=Win64+PS5+XSX -configuration=Development -build -skipcook -iterate -skippak -stage -server
