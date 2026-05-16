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
		RunDice(command[1:])
	default:
		fmt.Println("command not found: ", command[0])
		os.Exit(1)
	}
}
