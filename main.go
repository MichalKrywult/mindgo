package main

import "fmt"

func main() {
	response, err := NewMoodEntry(2, "2022/12/13", "Test")
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Println("Mood created correctly", response)
}
