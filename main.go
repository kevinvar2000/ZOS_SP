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

	for {
		printHelp()
		fmt.Print("Enter the command: ")
		var command string
		fmt.Scanln(&command)

		switch command {
		case "cpy":
			copyFile(input)
		case "mv":
			moveFile(input)
		case "rm":
			removeFile(input)
		case "mkdir":
			makeDirectory()
		case "rmdir":
			removeDirectory()
		case "ls":
			printDirectoryContents()
		case "cat":
			printFileContents(input)
		case "cd":
			changePath()
		case "pwd":
			printCurrentPath()
		case "info":
			printInformation()
		case "incp":
			incp()
		case "outcp":
			outcp()
		case "load":
			loadFile(input)
		case "format":
			formatFile(input)
		case "exit":
			return
		default:
			fmt.Println("Invalid command")
		}

	}

	// Copy the file
	// copyFile(input)

	// Move the file
	// moveFile(input)

	// Remove the file
	// removeFile(input)

	// Make a directory
	// makeDirectory()

	// Remove a directory
	// removeDirectory()

	// Print the contents of the directory
	// printDirectoryContents()

	// Print the contents of the file
	// printFileContents(input)

	// Change the path
	// changePath()

	// Print the current path
	// printCurrentPath()

	// Print the information
	// printInformation()

	// incp
	// incp()

	// outcp
	// outcp()

	// Load the file
	// loadFile(input)

	// Format the file
	// formatFile(input)

}

func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("cpy - Copy the file")
	fmt.Println("mv - Move the file")
	fmt.Println("rm - Remove the file")
	fmt.Println("mkdir - Make a directory")
	fmt.Println("rmdir - Remove a directory")
	fmt.Println("ls - Print the contents of the directory")
	fmt.Println("cat - Print the contents of the file")
	fmt.Println("cd - Change the path")
	fmt.Println("pwd - Print the current path")
	fmt.Println("info - Print the information")
	fmt.Println("incp - incp")
	fmt.Println("outcp - outcp")
	fmt.Println("load - Load the file")
	fmt.Println("format - Format the file")
	fmt.Println("exit - Exit the program")
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
