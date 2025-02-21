package main

// TODO: Create an asynchronous goroutine that creates a key value map of each executable and the path at which it is stored if the "PATH" env var is set.

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

// Check in PATH
var pathVal, pathIsSet = os.LookupEnv("PATH")
var pathCommands sync.Map

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
		if val, ok := pathCommands.Load(arg); ok {
			fmt.Printf("%s is %s\n", arg, val)
			return
		}

	}
	fmt.Printf("%s: not found\n", arg)
}

func ExecCommand(commandList []string) {
	cmd := exec.Command(commandList[0], commandList[1:]...)
	out, err := cmd.Output()
	output := string(out)
	output = strings.Replace(output, "\n", "", 0)
	if err == nil {
		// fmt.Println(string(out))
		fmt.Fprint(os.Stdout, output)
	}
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("EXITING EXEC COMMAND")
}

func VerboseCommand(commandList []string) {
	//fmt.Printf("Program was passed %d args (including program name).\n", len(commandList))
	fmt.Fprintf(os.Stdout, "Program was passed %d args (including program name).\n", len(commandList))
	for num, each := range commandList {
		if num == 0 {
			fmt.Fprintf(os.Stdout, "Arg #%d (program name): %s\n", num, each)
			continue
		}
		fmt.Fprintf(os.Stdout, "Arg #%d: %s\n", num, each)
	}
	// fmt.Println("EXITING VERBOSE COMMAND")
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
		if pathIsSet {
			// fmt.Println("[DEBUG]: switch hit default.")
			// check if command exists in path
			if _, ok := pathCommands.Load(commandList[0]); ok {
				// proceed to be verbose about input
				// VerboseCommand(commandList)
				ExecCommand(commandList)
				return
			}
		}
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
	// fmt.Println("[DEBUG]: in REPL")
	CheckCommand(message)
}

func main() {
	if pathIsSet {
		operatingSystem := runtime.GOOS
		var pathDirs []string
		switch operatingSystem {
		case "windows":
			pathDirs = strings.Split(pathVal, ";")
		case "linux":
			pathDirs = strings.Split(pathVal, ":")
		default:
			pathDirs = strings.Split(pathVal, ":")
		}
		var wg sync.WaitGroup
		// number of routines to spawn
		numRoutines := len(pathDirs)
		wg.Add(numRoutines)
		for v := range numRoutines {
			go func(v int) {
				defer wg.Done()
				//get source dir
				dir := pathDirs[v]
				files, err := os.ReadDir(dir)
				if err != nil {
					return
				}
				// Adds to the map in the following format:
				// "executableName":"path/to/executable"
				// Using Load or Store instead of just Store to overwrite to latest exe
				for _, file := range files {
					pathCommands.LoadOrStore(file.Name(), fmt.Sprintf("%s/%s", dir, file.Name()))

				}

			}(v)
		}
		wg.Wait()
	}

	for {
		REPL()
	}
}
