package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func ReadUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()
	message = strings.ReplaceAll(message, "\r\n", "")
	return message
}

func ListBuiltins(arg string) {
	supportedCommands := []string{"echo", "exit"}
	for _, v := range supportedCommands {
		if v == arg {
			fmt.Printf("%s is a shell builtin\n", v)
			return
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

func REPL() {
	fmt.Fprint(os.Stdout, "$ ")
	// Using scanner to read the input instead:
	message := ReadUserInput()
	CheckCommand(message)
}

func main() {
	for {
		REPL()
	}
}
