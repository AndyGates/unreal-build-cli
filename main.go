package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/multichoose"
)

func main() {

	config := GetConfig()
	uprojectPath := FindUproject()

	if uprojectPath == "" {
		log.Fatal("No .uproject found, make sure you are in the project directory")
	} else {
		log.Printf("Found project file: %s", uprojectPath)
	}

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

	// RunUAT.bat BuildCookRun -project="..."
	args := []string{"BuildCookRun"}

	// -project
	args = append(args, fmt.Sprintf("-project=\"%s\"", uprojectPath))

	// -targetplatform
	args = append(args, fmt.Sprintf("-targetplatform=%s", strings.Join(clientPlatforms, "+")))

	// -serverplatform
	if len(serverPlatforms) > 0 {
		args = append(args, "-server")
		args = append(args, fmt.Sprintf("-serverplatform=%s", strings.Join(serverPlatforms, "+")))
	}

	// -configuration
	args = append(args, fmt.Sprintf("-configuration=%s", strings.Join(configurations, "+")))

	// -cook -skipcook etc.
	for _, stepOption := range config.StepOptions.Options {
		stepArg := strings.ToLower(stepOption)
		if !Contains(steps, stepOption) {
			stepArg = fmt.Sprintf("skip%s", stepArg)
		}
		args = append(args, fmt.Sprintf("-%s", stepArg))
	}

	// -iterate
	if shouldIterate {
		args = append(args, "-iterate")
	}

	fmt.Println(args)

	cmd := exec.Command(config.RunUATPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
}

// ..\UnrealEngine\Engine\Build\BatchFiles\RunUAT.bat BuildCookRun -project="..\Yakisoba\Yakisoba.uproject" -platform=Win64+PS5+XSX -configuration=Development -build -skipcook -iterate -skippak -stage -server
