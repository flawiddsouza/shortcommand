package main

import (
	"bufio"
	"fmt"
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

func executeCommand(commandToRun string, currentWorkingDirectory string) error {
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	cmd := exec.Command("sh", "-c", commandToRun)
	if currentWorkingDirectory != "" {
		currentWorkingDirectoryResolved := os.ExpandEnv(currentWorkingDirectory)
		if strings.HasPrefix(currentWorkingDirectoryResolved, "~") {
			home, _ := os.UserHomeDir()
			currentWorkingDirectoryResolved = home + currentWorkingDirectoryResolved[1:]
		}
		if !fileExists(currentWorkingDirectoryResolved) {
			fmt.Printf("%sUnable to find given working directory: %s%s\n", colorRed, currentWorkingDirectory, colorReset)
		}
		cmd.Dir = currentWorkingDirectoryResolved
	}

	cmdReader, _ := cmd.StdoutPipe()

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s--> %s\n%s", colorGreen, scanner.Text(), colorReset)
		}
	}()

	cmdReader2, _ := cmd.StderrPipe()

	scanner2 := bufio.NewScanner(cmdReader2)
	go func() {
		for scanner2.Scan() {
			fmt.Printf("%s--> %s\n%s", colorRed, scanner2.Text(), colorReset)
		}
	}()

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
								fmt.Println("Not continuing with the next command in do sequence as previous command errored out")
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
