package main

// TODO: EchoFormatter needs to accept the updated arglist

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
var envHome, homeIsSet = os.LookupEnv("HOME")
var pathCommands sync.Map

func ExtractArgs(argString string) []string {
	var argsList []string
	var build strings.Builder
	toggleCapture := false
	for _, v := range argString {
		// Instead of block check (implemented below), check per char
		if v == ' ' && !toggleCapture {
			if build.Len() > 0 {
				argsList = append(argsList, build.String())
				build.Reset()
			}
			continue
		}
		if v == '\'' {
			toggleCapture = !toggleCapture
			continue
		}
		build.WriteRune(v)

	}
	if build.Len() > 0 {
		argsList = append(argsList, build.String())
		build.Reset()
	}

	return argsList
}

func EchoFormatter(printList []string) {
	fmt.Fprintln(os.Stdout, strings.Join(printList, " "))
}

func ReadUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()
	message = strings.ReplaceAll(message, "\r\n", "")
	return message
}

func ListBuiltins(arg string) {
	supportedCommands := []string{"echo", "exit", "type", "pwd", "cd"}

	for _, v := range supportedCommands {
		if v == arg {
			fmt.Fprintf(os.Stderr, "%s is a shell builtin\n", v)
			return
		}
	}
	if pathIsSet {
		if val, ok := pathCommands.Load(arg); ok {
			fmt.Fprintf(os.Stderr, "%s is %s\n", arg, val)
			return
		}

	}
	fmt.Fprintf(os.Stderr, "%s: not found\n", arg)
}

func ExecuteCommand(binaryPath string, commandList []string) {
	cmd := exec.Command(commandList[0], commandList[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func VerboseCommand(commandList []string) {
	//fmt.Printf("Program was passed %d args (including program name).\n", len(commandList))
	fmt.Fprintf(os.Stdout, "Program was passed %d args (including program name).\n", len(commandList))
	for num, each := range commandList {
		if num == 0 {
			if strings.ContainsRune(each, '/') || strings.ContainsRune(each, '\\') {
				var progName string
				var pathList []string
				operatingSystem := runtime.GOOS
				switch operatingSystem {
				case "windows":
					pathList = strings.SplitAfterN(each, "\\", 1)
				case "linux":
					pathList = strings.Split(each, "/")
				default:
					pathList = strings.Split(each, "/")
				}
				progName = pathList[len(pathList)-1]
				fmt.Fprintf(os.Stdout, "Arg #0 (program name): %s\n", progName)
				continue
			}
			fmt.Fprintf(os.Stdout, "Arg #0 (program name): %s\n", commandList[0])
			continue
		}
		fmt.Fprintf(os.Stdout, "Arg #%d: %s\n", num, each)
	}
}

func CheckCommand(command string) {
	commandList := strings.Split(command, " ")
	executable := commandList[0]
	args := ExtractArgs(strings.Join(commandList[1:], " "))
	commandList = nil
	commandList = append(commandList, executable)
	commandList = append(commandList, args...)
	switch commandList[0] {
	case "exit":
		code, err := strconv.ParseInt(commandList[1], 10, 16)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(0)
		}
		os.Exit(int(code))
	case "echo":
		EchoFormatter(commandList[1:])
	case "type":
		ListBuiltins(commandList[1])
	case "pwd":
		if pwd, ok := os.Getwd(); ok == nil {
			fmt.Fprintln(os.Stdout, pwd)
		}
	case "cd":
		if len(commandList[1:]) == 1 {
			if commandList[1] == "~" {
				if homeIsSet {
					err := os.Chdir(envHome)
					if err != nil {
						customErr := strings.Split(err.Error(), " ")
						// This is a little inappropriate but conforming to codecrafters for now
						customErr[0] = "cd:"
						customErr[2] = "No"
						outString := strings.Join(customErr, " ")
						fmt.Fprintln(os.Stdout, outString)
						return
					}
					return
				}
			}
			err := os.Chdir(commandList[1])
			if err != nil {
				customErr := strings.Split(err.Error(), " ")
				// This is a little inappropriate but conforming to codecrafters for now
				customErr[0] = "cd:"
				customErr[2] = "No"
				outString := strings.Join(customErr, " ")
				fmt.Fprintln(os.Stdout, outString)
				return
			}
		}
	case "cat":
		ExecuteCommand("cat", commandList)
	default:
		if pathIsSet {
			// check if command exists in path
			if binaryPath, ok := pathCommands.Load(commandList[0]); ok {
				// proceed to be verbose about input
				// NOTE: commented out due to an interesting bug in codecrafters tester.
				//VerboseCommand(commandList)
				ExecuteCommand(binaryPath.(string), commandList)
				return
			}
		}
		_, err := fmt.Printf("%s: command not found\n", commandList[0])
		if err != nil {
			fmt.Fprint(os.Stderr, err)
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
				//traverse the $PATH var in reverse to ensure the latest (traditionally) binary is the only one in list
				dir := pathDirs[len(pathDirs)-v-1]
				files, err := os.ReadDir(dir)
				if err != nil {
					return
				}
				// Adds to the map in the following format:
				// "executableName":"path/to/executable"
				// Using Load or Store instead of just Store to overwrite to latest exe
				for _, file := range files {
					// pathCommands.LoadOrStore(file.Name(), fmt.Sprintf("%s/%s", dir, file.Name()))
					// Choosing store over load or store since we want the overwrite.
					pathCommands.Store(file.Name(), fmt.Sprintf("%s/%s", dir, file.Name()))
				}

			}(v)
		}
		wg.Wait()
	}

	for {
		REPL()
	}
}
