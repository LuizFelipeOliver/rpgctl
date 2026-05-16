package main

import (
	"fmt"
	"os"
)

func main() {
	command := os.Args[1:]
	if len(command) == 0 {
		fmt.Println("rpgctl <command> [arguments]")
		os.Exit(2)
	}
	switch command[0] {
	case "dice":
		if err := RunDice(command[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Println("command not found: ", command[0])
		os.Exit(1)
	}
}
