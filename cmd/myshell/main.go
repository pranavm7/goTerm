package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	reader:= bufio.NewReader(os.Stdin)
	message,_:=  reader.ReadString('\n')
	message= strings.ReplaceAll(message,"\r\n","")
	// fmt.Fprintln(os.Stdout, name, "is", age, "years old.")
	fmt.Fprintf(os.Stdout,  "%s: Command not found\n",message)
}
