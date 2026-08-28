package main

import "os"

func main() {
	tracker := MoodTracker{}
	cli := NewCLI(&tracker, os.Stdin)
	cli.show() // we have to use &, so CLI function will operate on the same MoodTracker object - not a copy of it
}
