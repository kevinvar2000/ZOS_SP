package main

import (
	"fmt"
)

func main() {

	// Read arguments from the command line
	// args := os.Args
	// if len(args) < 2 {
	// 	fmt.Println("Usage: go run main.go <file_name>")
	// 	return
	// }

	// Read input from the user
	var input string
	fmt.Print("Enter the file name: ")
	fmt.Scanln(&input)

	// Copy the file
	copyFile(input)

	// Move the file
	moveFile(input)

	// Remove the file
	removeFile(input)

	// Make a directory
	makeDirectory()

	// Remove a directory
	removeDirectory()

	// Print the contents of the directory
	printDirectoryContents()

	// Print the contents of the file
	printFileContents(input)

	// Change the path
	changePath()

	// Print the current path
	printCurrentPath()

	// Print the information
	printInformation()

	// incp
	incp()

	// outcp
	outcp()

	// Load the file
	loadFile(input)

	// Format the file
	formatFile(input)

}

func formatFile(input string) {
	panic("unimplemented")
}

func loadFile(input string) {
	panic("unimplemented")
}

func outcp() {
	panic("unimplemented")
}

func incp() {
	panic("unimplemented")
}

func printInformation() {
	panic("unimplemented")
}

func printCurrentPath() {
	panic("unimplemented")
}

func changePath() {
	panic("unimplemented")
}

func printFileContents(input string) {
	panic("unimplemented")
}

func printDirectoryContents() {
	panic("unimplemented")
}

func removeDirectory() {
	panic("unimplemented")
}

func makeDirectory() {
	panic("unimplemented")
}

func removeFile(input string) {
	panic("unimplemented")
}

func moveFile(input string) {
	panic("unimplemented")
}

func copyFile(input string) {
	panic("unimplemented")
}
