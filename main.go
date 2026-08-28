package main

import (
	"bufio"
	"os"
)

func main() {
	tracker := MoodTracker{}
	scanner := bufio.NewScanner(os.Stdin)

	showCLI(&tracker, scanner) // we have to use &, so CLI function will operate on the same MoodTracker object - not a copy of it
}
