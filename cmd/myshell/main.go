package main

// TODO: Create an asynchronous goroutine that creates a key value map of each executable and the path at which it is stored if the "PATH" env var is set.

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

// Check in PATH
var pathVal, pathIsSet = os.LookupEnv("PATH")

func ReadUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()
	message = strings.ReplaceAll(message, "\r\n", "")
	return message
}

func ListBuiltins(arg string) {
	supportedCommands := []string{"echo", "exit", "type"}

	for _, v := range supportedCommands {
		if v == arg {
			fmt.Printf("%s is a shell builtin\n", v)
			return
		}
	}
	if pathIsSet {
		// TODO: Write a lookup to the execCache
		// Get a list of path directories.
		pathDirs := strings.Split(pathVal, ":")
		for _, pathDirectory := range pathDirs {
			// TODO: Iterate over paths
		}
	}
	fmt.Printf("%s: not found\n", arg)
}

func CheckCommand(command string) {
	commandList := strings.Split(command, " ")
	// Check for exit
	switch commandList[0] {
	case "exit":
		code, err := strconv.ParseInt(commandList[1], 10, 16)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		os.Exit(int(code))
	case "echo":
		fmt.Println(strings.Join(commandList[1:], " "))
	case "type":
		ListBuiltins(commandList[1])
	default:
		_, err := fmt.Printf("%s: command not found\n", commandList[0])
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}
}

// building the READ EVAL PRINT LOOP
func REPL() {
	fmt.Fprint(os.Stdout, "$ ")
	message := ReadUserInput()
	CheckCommand(message)
}

func main() {
	if pathIsSet {
		// TODO: Call the goroutines to list all the executables and add to a mutex

	}
	for {
		REPL()
	}
}
