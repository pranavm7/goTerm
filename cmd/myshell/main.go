package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	//"golang.org/x/text/message"
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

func CheckCommand(command string) {
	// Check for exit
	switch command {
	case "exit 0":
		os.Exit(0)
	}
}

func REPL() {
	fmt.Fprint(os.Stdout, "$ ")
	// Using scanner to read the input instead:
	message := ReadUserInput()
	CheckCommand(message)
	message = fmt.Sprintf("%s: command not found", message)
	_, _ = fmt.Println(message)
}

func main() {
	// Uncomment this block to pass the first stage
	// fmt.Fprint(os.Stdout, "$ ")

	// // Using scanner to read the input instead:
	// message:= ReadUserInput()
	// message= fmt.Sprintf("%s: command not found",message)
	// _,_ =fmt.Println(message)
	for {
		REPL()
	}
}
