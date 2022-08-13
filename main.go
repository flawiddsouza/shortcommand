package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ShortCommands []ShortCommand `yaml:"shortcommands"`
}

type ShortCommand struct {
	Name     string    `yaml:"name"`
	Commands []Command `yaml:"commands"`
}

type Command struct {
	Name                    string   `yaml:"name"`
	Description             string   `yaml:"description"`
	CurrentWorkingDirectory string   `yaml:"cwd"`
	Do                      []string `yaml:"do"`
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

var ConsoleColors = map[string]string{
	"Reset": "\033[0m",
	"Red":   "\033[31m",
}

func executeCommand(commandToRun string, currentWorkingDirectory string) error {
	cmd := exec.Command("sh", "-c", commandToRun)
	if currentWorkingDirectory != "" {
		currentWorkingDirectoryResolved := os.ExpandEnv(currentWorkingDirectory)
		if strings.HasPrefix(currentWorkingDirectoryResolved, "~") {
			home, _ := os.UserHomeDir()
			currentWorkingDirectoryResolved = home + currentWorkingDirectoryResolved[1:]
		}
		if !fileExists(currentWorkingDirectoryResolved) {
			fmt.Printf("%sUnable to find given working directory: %s%s\n", ConsoleColors["Red"], currentWorkingDirectory, ConsoleColors["Reset"])
		}
		cmd.Dir = currentWorkingDirectoryResolved
	}

	stdout := io.Writer(os.Stdout)
	stderr := io.Writer(os.Stderr)

	cmd.Stderr = stderr
	cmd.Stdout = stdout

	cmd.Start()

	return cmd.Wait()
}

func main() {
	shortcommandPath := os.Getenv("SHORTCOMMAND_CONFIG")

	if shortcommandPath == "" {
		fmt.Println("SHORTCOMMAND_CONFIG env not found or empty")
		return
	}

	yamlFile, err := os.ReadFile(shortcommandPath)

	if err != nil {
		fmt.Println("Unable to find file specified in SHORTCOMMAND_CONFIG: ", shortcommandPath)
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Two arguments need to passed to the command, first is your shortcommand and then the command under it")
		return
	}

	var parsedConfig Config

	err = yaml.Unmarshal(yamlFile, &parsedConfig)
	if err != nil {
		fmt.Println("Invalid config file given")
		return
	}

	for _, shortCommand := range parsedConfig.ShortCommands {
		if shortCommand.Name == os.Args[1] {
			for _, command := range shortCommand.Commands {
				if command.Name == os.Args[2] {
					fmt.Println(shortCommand.Name, ">", command.Name, ":", command.Description)
					for i, commandToRun := range command.Do {
						fmt.Println(">", commandToRun)
						err = executeCommand(commandToRun, command.CurrentWorkingDirectory)
						if err != nil {
							if i < len(command.Do)-1 {
								fmt.Printf("%sNot continuing with the next command in do sequence as previous command errored out%s\n", ConsoleColors["Red"], ConsoleColors["Reset"])
								return
							}
						}
					}
					fmt.Println()
					return
				}
			}
			break
		}
	}

	fmt.Println("No matching command found for given arguments")
}
