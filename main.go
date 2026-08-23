package main

func main() {
	tracker := MoodTracker{}
	showCLI(&tracker) // we have to use &, so CLI function will operate on the same MoodTracker object - not a copy of it
}
