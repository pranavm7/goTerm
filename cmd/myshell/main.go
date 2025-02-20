package main

// TODO: Create an asynchronous goroutine that creates a key value map of each executable and the path at which it is stored if the "PATH" env var is set.

import (
	"bufio"
	"fmt"
	"os"
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
		// TODO: Write a lookup to the execCache
		if val, ok := pathCommands.Load(arg); ok {
			fmt.Printf("%s is %s", arg, val)
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

// building the READ EVAL PRINT LOOP
func REPL() {
	fmt.Fprint(os.Stdout, "$ ")
	message := ReadUserInput()
	CheckCommand(message)
}

func main() {
	if pathIsSet {
		operatingSystem := runtime.GOOS
		pathDirs := []string{}
		switch operatingSystem {
		case "windows":
			pathDirs = strings.Split(pathVal, ";")
		case "linux":
			pathDirs = strings.Split(pathVal, ":")
		default:
			pathDirs = strings.Split(pathVal, ":")
		}
		// pathDirs := strings.Split(pathVal, ":")
		// TODO: Call the goroutines to list all the executables and add to a mutex
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
					pathCommands.LoadOrStore(file.Name(), fmt.Sprintf("%s/%s\n", dir, file.Name()))

				}

			}(v)
		}
		wg.Wait()
		// printing output
		// m := map[string]interface{}{}
		// pathCommands.Range(func(key, value interface{}) bool {
		// 	m[fmt.Sprint(key)] = value
		// 	return true
		// })

		// b, err := json.MarshalIndent(m, "", " ")
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(string(b))
		// fmt.Println()
		// fmt.Println(pathDirs)
		// fmt.Println()
		// fmt.Println(pathVal)
	}

	for {
		REPL()
	}
}
