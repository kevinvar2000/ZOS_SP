package main

import (
	"fmt"
	"os"
	"strings"
)

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

func formatFile(filename string) {
	// Reset the FAT table and directory
	fs := &FileSystem{}
	fs.Init()

	// Optionally, you can write the empty FAT and directory to disk
	// to persist the formatted file system.
	fmt.Println("OK")
}

func loadFile(filename string) {
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

func outcp(src string, dest string) {
	entry, exists := fs.Directory[src]
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
	currentCluster := entry.FirstCluster
	for currentCluster != FAT_EOF {
		outFile.Write(fs.ClusterData[currentCluster])
		currentCluster = fs.FatTable[currentCluster]
	}

	fmt.Println("OK")
}

func incp(src string, dest string) {
	// Read the source file from the host system
	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// Copy the file into the pseudoFAT system
	err = fs.inCp(dest, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("OK")
}

func printInformation(filename string) {
	entry, exists := fs.Directory[filename]
	if !exists {
		fmt.Println("FILE NOT FOUND")
		return
	}

	fmt.Printf("%s:", filename)

	// Traverse the FAT chain and print cluster numbers
	currentCluster := entry.FirstCluster
	for currentCluster != FAT_EOF {
		fmt.Printf(" %d", currentCluster)
		currentCluster = fs.FatTable[currentCluster]
	}
	fmt.Println()
}

func printCurrentPath() {
	// Assume a simple path system, or modify it based on your directory structure
	fmt.Println(currentPath) // currentPath should be a global variable maintaining the current path
}

func changePath(newPath string) {
	// Check if the directory exists in the pseudoFAT system
	_, exists := fs.Directory[newPath]
	if !exists {
		fmt.Println("PATH NOT FOUND")
		return
	}

	// Update currentPath to the new path
	currentPath = newPath
	fmt.Println("OK")
}

func printFileContents(filename string) {
	err := fs.cat(filename)
	if err != nil {
		fmt.Println(err)
	}
}

func printDirectoryContents() {
	for name, entry := range fs.Directory {
		fmt.Printf("FILE: %s, SIZE: %d bytes\n", name, entry.Size)
	}
}

func removeDirectory(dirname string) {
	// Check if the directory is empty
	for name := range fs.Directory {
		if strings.HasPrefix(name, dirname) {
			fmt.Println("NOT EMPTY")
			return
		}
	}

	// Remove the directory from the directory list
	delete(fs.Directory, dirname)
	fmt.Println("OK")
}

func makeDirectory(dirname string) {
	// Ensure directory doesn't already exist
	if _, exists := fs.Directory[dirname]; exists {
		fmt.Println("EXIST")
		return
	}

	// Create a directory (it behaves as a special file)
	fs.Directory[dirname] = DirectoryEntry{Name: dirname, Size: 0, FirstCluster: FAT_EOF}
	fmt.Println("OK")
}

func removeFile(filename string) {
	err := fs.rm(filename)
	if err != nil {
		fmt.Println(err)
	}
}

func moveFile(src string, dest string) {
	// Ensure the source file exists
	entry, exists := fs.Directory[src]
	if !exists {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// Ensure the destination does not exist
	if _, exists := fs.Directory[dest]; exists {
		fmt.Println("PATH NOT FOUND")
		return
	}

	// Rename the file (move in directory)
	fs.Directory[dest] = entry
	delete(fs.Directory, src)
	fmt.Println("OK")
}

func copyFile(src string, dest string) {
	// Ensure the source file exists
	entry, exists := fs.Directory[src]
	if !exists {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// Ensure the destination does not exist
	if _, exists := fs.Directory[dest]; exists {
		fmt.Println("PATH NOT FOUND")
		return
	}

	// Read the file data from source
	currentCluster := entry.FirstCluster
	var data []byte
	for currentCluster != FAT_EOF {
		data = append(data, fs.ClusterData[currentCluster]...)
		currentCluster = fs.FatTable[currentCluster]
	}

	// Copy the file data to the destination
	err := fs.inCp(dest, data)
	if err != nil {
		fmt.Println(err)
	}
}
