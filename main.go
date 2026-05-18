package main

import (
	"fmt"
	"os"
)

func main() {
	command := os.Args[1:]
	if len(command) == 0 {
		fmt.Println("Usage: rpgctl <command> [arguments]")
		fmt.Println("Available commands:")
		fmt.Println("  dice  - Roll dice using NdM notation")
		os.Exit(2)
	}
	switch command[0] {
	case "dice":
		if err := RunDice(command[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "init":
		if err := RunInitiative(command[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "command not found:", command[0])
		os.Exit(1)
	}
}
