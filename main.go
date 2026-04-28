package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	r  = color.New(color.FgRed)
	g  = color.New(color.FgGreen)
	gb = color.New(color.FgGreen, color.Bold)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("\tNot enough args!\n\tIf you wish to capture data, do:\n\t\t%s capture\n", os.Args[0])
		fmt.Printf("\tIf you wish to analyze recorded data, do:\n\t\t%s <source_dir>\n", os.Args[0])
		return
	}
	if strings.EqualFold(os.Args[1], "capture") {
		err := BeginCapture()
		if err != nil {
			r.Print("\t[ERROR] ")
			fmt.Printf("%v\n", err)
		}
		return
	}

	session, err := LoadSession(os.Args[1])
	if err != nil {
		r.Print("\t[ERROR] ")
		fmt.Printf("Failed to load data: %v\n", err)
		return
	}

	// begin cli loop
	for {
		gb.Print(" $ ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}
		tokens := strings.Fields(command)
		exit, err := ParseCommand(tokens)
		if err != nil {
			fmt.Printf("Invalid command: %v\n", err)
		}
		if exit {
			break
		}
	}
}
