package main

import (
	"fmt"
	"os"
	"strings"
)

func PrintHelp() {
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

func FormatFileCmd(filename string) {
	// Reset the FAT table and directory
	fs := &FileSystem{}
	fs.Init()

	// File doesn't exist, ask for the file size and create it
	var fileSize int64 // TODO: Change to bytes
	fmt.Print("Enter the desired file size in bytes: ")
	fmt.Scanln(&fileSize)

	FormatFile(filename, fileSize)

	// Optionally, you can write the empty FAT and directory to disk
	// to persist the formatted file system.
	fmt.Println("OK")
}

func LoadFile(filename string) {
	// Open the file that contains commands
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("FILE NOT FOUND")
		return
	}

	commands := strings.Split(string(data), "\n")
	for _, command := range commands {
		// Execute each command
		// Example: handle cp, mv, etc. by splitting the command and calling the appropriate function
		fmt.Println("Executing:", command)
	}

	fmt.Println("OK")
}

func Outcp(src string, dest string) {
	entry, exists := fs.directory[src]
	if !exists {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// Open the destination file on the host file system
	outFile, err := os.Create(dest)
	if err != nil {
		fmt.Println("PATH NOT FOUND")
		return
	}
	defer outFile.Close()

	// Read from the pseudoFAT and write to the destination file
	currentCluster := entry.first_cluster
	for currentCluster != FAT_EOF {
		outFile.Write(fs.cluster_data[currentCluster])
		currentCluster = fs.fat_table[currentCluster]
	}

	fmt.Println("OK")
}

func Incp(src string, dest string) {
	// Read the source file from the host system
	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// Copy the file into the pseudoFAT system
	err = fs.InCp(dest, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("OK")
}

func PrintInformation(filename string) {
	entry, exists := fs.directory[filename]
	if !exists {
		fmt.Println("FILE NOT FOUND")
		return
	}

	fmt.Printf("%s:", filename)

	// Traverse the FAT chain and print cluster numbers
	currentCluster := entry.first_cluster
	for currentCluster != FAT_EOF {
		fmt.Printf(" %d", currentCluster)
		currentCluster = fs.fat_table[currentCluster]
	}
	fmt.Println()
}

func PrintCurrentPath() {
	// Assume a simple path system, or modify it based on your directory structure
	fmt.Println(currentPath) // currentPath should be a global variable maintaining the current path
}

func ChangePath(newPath string) {
	// Check if the directory exists in the pseudoFAT system
	_, exists := fs.directory[newPath]
	if !exists {
		fmt.Println("PATH NOT FOUND")
		return
	}

	// Update currentPath to the new path
	currentPath = newPath
	fmt.Println("OK")
}

func PrintFileContents(filename string) {
	err := fs.Cat(filename)
	if err != nil {
		fmt.Println(err)
	}
}

func PrintDirectoryContents() {
	for name, entry := range fs.directory {
		fmt.Printf("FILE: %s, SIZE: %d bytes\n", name, entry.size)
	}
}

func RemoveDirectory(dirname string) {
	// Check if the directory is empty
	for name := range fs.directory {
		if strings.HasPrefix(name, dirname) {
			fmt.Println("NOT EMPTY")
			return
		}
	}

	// Remove the directory from the directory list
	delete(fs.directory, dirname)
	fmt.Println("OK")
}

func MakeDirectory(dirname string) {
	// Ensure directory doesn't already exist
	if _, exists := fs.directory[dirname]; exists {
		fmt.Println("EXIST")
		return
	}

	// Create a directory (it behaves as a special file)
	fs.directory[dirname] = DirectoryEntry{name: dirname, size: 0, first_cluster: FAT_EOF}
	fmt.Println("OK")
}

func RemoveFile(filename string) {
	err := fs.Rm(filename)
	if err != nil {
		fmt.Println(err)
	}
}

func MoveFile(src string, dest string) {
	// Ensure the source file exists
	entry, exists := fs.directory[src]
	if !exists {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// Ensure the destination does not exist
	if _, exists := fs.directory[dest]; exists {
		fmt.Println("PATH NOT FOUND")
		return
	}

	// Rename the file (move in directory)
	fs.directory[dest] = entry
	delete(fs.directory, src)
	fmt.Println("OK")
}

func CopyFile(src string, dest string) {
	// Ensure the source file exists
	entry, exists := fs.directory[src]
	if !exists {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// Ensure the destination does not exist
	if _, exists := fs.directory[dest]; exists {
		fmt.Println("PATH NOT FOUND")
		return
	}

	// Read the file data from source
	currentCluster := entry.first_cluster
	var data []byte
	for currentCluster != FAT_EOF {
		data = append(data, fs.cluster_data[currentCluster]...)
		currentCluster = fs.fat_table[currentCluster]
	}

	// Copy the file data to the destination
	err := fs.InCp(dest, data)
	if err != nil {
		fmt.Println(err)
	}
}
