package main

import (
	"fmt"
	"os"
)

func main() {

	storage := FileStorage{filename: "moods.json"}
	tracker, err := NewMoodTracker(storage)
	if err != nil {
		fmt.Println("Error initializing tracker:", err)
		return
	}

	cli := NewCLI(tracker, os.Stdin)// the tracker is a pointer, received from NewMoodTracker function
	cli.show()
}
