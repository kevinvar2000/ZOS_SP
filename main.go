package main

import (
	"fmt"
	"os"
	"strings"
)

func enterCommand() {

	for {
		PrintHelp()
		fmt.Print("Enter the command: ")
		var command string
		fmt.Scanln(&command)

		switch command {
		case "cpy":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			CopyFile(src, dest)
		case "mv":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			MoveFile(src, dest)
		case "rm":
			fmt.Print("Enter the file to remove: ")
			var filename string
			fmt.Scanln(&filename)
			RemoveFile(filename)
		case "mkdir":
			fmt.Print("Enter the directory name: ")
			var dirname string
			fmt.Scanln(&dirname)
			MakeDirectory(dirname)
		case "rmdir":
			fmt.Print("Enter the directory name: ")
			var dirname string
			fmt.Scanln(&dirname)
			RemoveDirectory(dirname)
		case "ls":
			PrintDirectoryContents()
		case "cat":
			fmt.Print("Enter the file to display: ")
			var filename string
			fmt.Scanln(&filename)
			PrintFileContents(filename)
		case "cd":
			fmt.Print("Enter the new path: ")
			var newPath string
			fmt.Scanln(&newPath)
			ChangePath(newPath)
		case "pwd":
			PrintCurrentPath()
		case "info":
			fmt.Print("Enter the file to get info: ")
			var filename string
			fmt.Scanln(&filename)
			PrintInformation(filename)
		case "incp":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			Incp(src, dest)
		case "outcp":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			Outcp(src, dest)
		case "load":
			fmt.Print("Enter the script filename to load: ")
			var scriptFilename string
			fmt.Scanln(&scriptFilename)
			LoadFile(scriptFilename)
		case "format":
			fmt.Print("Enter the filesystem filename to format: ")
			var fsFilename string
			fmt.Scanln(&fsFilename)
			FormatFileCmd(fsFilename)
		case "exit":
			fmt.Println("Exiting the file system simulator.")
			return
		default:
			fmt.Println("Invalid command")
		}

	}

}

func checkFilename() string {

	var filename string

	// Loop until a valid filename with the ".dat" extension is provided
	for {
		if filename == "" {
			// Prompt the user to enter the filename
			fmt.Print("Enter the file name: ")
			fmt.Scanln(&filename)
		}

		// Check if the file has the correct ".dat" extension
		if strings.HasSuffix(filename, ".dat") {
			return filename
		}

		// Invalid extension, re-prompt the user
		fmt.Println("Invalid file extension. Please use a .dat file.")
		filename = ""
	}

}

func checkFile(filename string) {

	info, err := os.Stat(filename)
	if err == nil {

		// Check if the file is empty
		if info.Size() == 0 {
			fmt.Println("File is empty.")
		} else {
			// Load the file system
			LoadFileSystem(filename)

		}

	} else if os.IsNotExist(err) {

		fmt.Printf("\nFile does not exist. Formatting a new file system...\n")

		// Format the file with the desired size
		Format(filename)

		fmt.Printf("File created successfully!\n\n")

	} else {

		fmt.Println("Error reading file:", err)

	}

}

func main() {

	fmt.Printf("Welcome to the file system simulator\n")
	fmt.Printf("KIV/ZOS - SP 2024; Author: Kevin Varchola\n\n")

	// Read arguments from the command line
	args := os.Args
	var filename string

	// Check if the filename is provided as a command-line argument
	if len(args) == 2 {
		filename = args[1]
	} else {
		fmt.Println("Usage: go run main.go <file_name>")
		fmt.Println("Please provide a file name as an argument.")
		filename = checkFilename()
		// return
	}

	checkFile(filename)

	// enterCommand()

}
