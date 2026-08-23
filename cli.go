package main

import "fmt"

func showCLI() {
	fmt.Println("=====MENU=====")
	fmt.Println("1. New entry")
	fmt.Println("2. Remove entry")
	fmt.Println("3. Show history")
	fmt.Println("4. Show statistics")

	fmt.Println("Your choice: ")
	var choice int
	fmt.Scan(&choice) //& means using adress of a variable, Scan saves input to that variable

	switch choice {
	case 1:
		fmt.Println("Choice 1")
	case 2:
		fmt.Println("Choice 2")
	case 3:
		fmt.Println("Choice 3")
	case 4:
		fmt.Println("Choice 4")
	default:
		fmt.Println("Invalid choice")
	}
}
