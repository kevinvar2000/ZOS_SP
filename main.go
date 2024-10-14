package main

import (
	"fmt"
)

var currentPath string = "/" // Assume a simple path system
var fs *FileSystem           // Declare globally to use in all functions

func main() {

	// Read arguments from the command line
	// args := os.Args
	fs := &FileSystem{}
	fs.Init()

	// Read arguments from the command line
	// 	fmt.Println("Usage: go run main.go <file_name>")
	// 	return
	// }

	fmt.Println("Welcome to the file system simulator")
	fmt.Println("KIV/ZOS - SP 2024; Author: Kevin Varchola")

	// Read filename from the user
	// var filename string
	// fmt.Print("Enter the file name: ")
	// fmt.Scanln(&filename)

	// _, err := os.Stat(filename)
	// if err == nil {

	// 	// File exists, read its contents
	// 	data, err := os.ReadFile(filename)
	// 	if err != nil {
	// 		fmt.Println("Error reading file:", err)
	// 		return
	// 	}
	// 	fmt.Println("File name:", filename)
	// 	fmt.Printf("File size: %d bytes\n", len(data))
	// 	fmt.Printf("File content:\n%s\n", data)

	// } else if os.IsNotExist(err) {

	// 	// File doesn't exist, ask for the file size and create it
	// 	var fileSize int64
	// 	fmt.Print("Enter the desired file size in bytes: ")
	// 	fmt.Scanln(&fileSize)

	// 	file, err := os.Create(filename)
	// 	if err != nil {
	// 		fmt.Println("Error creating file:", err)
	// 		return
	// 	}
	// 	defer file.Close()

	// 	// Write the desired number of bytes to the file
	// 	if fileSize > 0 {
	// 		_, err = file.Write(make([]byte, fileSize))
	// 		if err != nil {
	// 			fmt.Println("Error writing to file:", err)
	// 			return
	// 		}
	// 	}

	// 	fmt.Println("File created successfully!")
	// }

	for {
		printHelp()
		fmt.Print("Enter the command: ")
		var command string
		fmt.Scanln(&command)

		switch command {
		case "cpy":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			copyFile(src, dest)
		case "mv":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			moveFile(src, dest)
		case "rm":
			fmt.Print("Enter the file to remove: ")
			var filename string
			fmt.Scanln(&filename)
			removeFile(filename)
		case "mkdir":
			fmt.Print("Enter the directory name: ")
			var dirname string
			fmt.Scanln(&dirname)
			makeDirectory(dirname)
		case "rmdir":
			fmt.Print("Enter the directory name: ")
			var dirname string
			fmt.Scanln(&dirname)
			removeDirectory(dirname)
		case "ls":
			printDirectoryContents()
		case "cat":
			fmt.Print("Enter the file to display: ")
			var filename string
			fmt.Scanln(&filename)
			printFileContents(filename)
		case "cd":
			fmt.Print("Enter the new path: ")
			var newPath string
			fmt.Scanln(&newPath)
			changePath(newPath)
		case "pwd":
			printCurrentPath()
		case "info":
			fmt.Print("Enter the file to get info: ")
			var filename string
			fmt.Scanln(&filename)
			printInformation(filename)
		case "incp":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			incp(src, dest)
		case "outcp":
			fmt.Print("Enter source and destination: ")
			var src, dest string
			fmt.Scanln(&src, &dest)
			outcp(src, dest)
		case "load":
			fmt.Print("Enter the script filename to load: ")
			var scriptFilename string
			fmt.Scanln(&scriptFilename)
			loadFile(scriptFilename)
		case "format":
			fmt.Print("Enter the filesystem filename to format: ")
			var fsFilename string
			fmt.Scanln(&fsFilename)
			formatFile(fsFilename)
		case "exit":
			fmt.Println("Exiting the file system simulator.")
			return
		default:
			fmt.Println("Invalid command")
		}

	}

}
