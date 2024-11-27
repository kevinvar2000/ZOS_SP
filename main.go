package main

import (
	"fmt"
	"os"
	"strings"
)

func enterCommand(filename string, fs_format FileSystemFormat) {

	PrintHelp()

	for {
		fmt.Print("Enter the command: ")
		var command, arg1, arg2 string
		fmt.Scanln(&command, &arg1, &arg2)
		ExecuteCommand(filename, command, arg1, arg2, fs_format)
		if command == "exit" || command == "quit" || command == "q" {
			break
		}
	}
}

func checkFilename() string {

	var filename string

	for {
		if filename == "" {
			fmt.Print("Enter the file name: ")
			fmt.Scanln(&filename)
		}

		if strings.HasSuffix(filename, ".dat") {
			return filename
		}

		fmt.Println("Invalid file extension. Please use a .dat file.")
		filename = ""
	}

}

func checkFile(filename string) {

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {

		fmt.Printf("\nFile does not exist. Formatting a new file system...\n")

		var file_size_mb int
		fmt.Print("Enter the desired file size in MB: ")
		fmt.Scanln(&file_size_mb)

		Format(filename, file_size_mb)
		fmt.Printf("File created and formatted successfully.\n\n")

	} else if err != nil {
		fmt.Println("Error reading file:", err)
	}

}

func main() {

	fmt.Printf("Welcome to the file system simulator\n")
	fmt.Printf("KIV/ZOS - SP 2024; Author: Kevin Varchola\n\n")

	args := os.Args
	var filename string

	if len(args) == 2 {
		filename = args[1]
	} else {
		fmt.Println("Usage: go run main.go <file_name>")
		fmt.Println("Please provide a file name as an argument.")
		filename = checkFilename()
		// return
	}

	checkFile(filename)

	fs_format := LoadFormat(filename)

	SetCurrentCluster(fs_format.data_start / CLUSTER_SIZE)

	enterCommand(filename, fs_format)

	fat1, fat2 := LoadFileSystem(filename)
	PrintFileSystem(fat1, fat2, "fats_after.txt")

}
