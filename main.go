package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Constants for file system
const (
	ClusterSize   = 512  // Each cluster is 512 bytes
	MaxFileName   = 12   // 8.3 format = 11 chars + null terminator
	MaxClusters   = 1024 // Max number of clusters in the file system
	FAT_FREE      = -1   // FAT free cluster marker
	FAT_EOF       = -2   // FAT end of file marker
	MaxFileLength = 5000 // Maximum file length for the `short` command
)

// FAT entry struct to simulate FAT table
type FAT [MaxClusters]int

// DirectoryEntry stores file metadata
type DirectoryEntry struct {
	Name         string
	Size         int
	FirstCluster int
	IsDirectory  bool
}

// FileSystem struct
type FileSystem struct {
	FatTable    FAT
	Directory   map[string]DirectoryEntry
	ClusterData [][]byte // Storage for the actual file data
}

// Initializes the file system
func (fs *FileSystem) Init() {

	fs.FatTable = [MaxClusters]int{}

	for i := range fs.FatTable {
		fs.FatTable[i] = FAT_FREE
	}

	fs.Directory = make(map[string]DirectoryEntry)
	fs.ClusterData = make([][]byte, MaxClusters)

	for i := range fs.ClusterData {
		fs.ClusterData[i] = make([]byte, ClusterSize)
	}

}

// Find a free cluster in FAT
func (fs *FileSystem) findFreeCluster() (int, error) {

	for i, val := range fs.FatTable {
		if val == FAT_FREE {
			return i, nil
		}
	}

	return -1, errors.New("no free clusters available")
}

// Load file from external system into pseudo-FAT (incp s1 s2)
func (fs *FileSystem) inCp(filename string, data []byte) error {

	// Ensure the filename is not too long
	if len(filename) > MaxFileName {
		return fmt.Errorf("filename too long")
	}

	// Ensure the file does not already exist
	if _, exists := fs.Directory[filename]; exists {
		return fmt.Errorf("file already exists")
	}

	// Split data into clusters
	numClusters := (len(data) + ClusterSize - 1) / ClusterSize
	firstCluster, err := fs.findFreeCluster()
	if err != nil {
		return err
	}

	currentCluster := firstCluster
	for i := 0; i < numClusters; i++ {
		if i > 0 {
			newCluster, err := fs.findFreeCluster()
			if err != nil {
				return err
			}
			fs.FatTable[currentCluster] = newCluster
			currentCluster = newCluster
		}

		// Copy the data into the cluster
		start := i * ClusterSize
		end := start + ClusterSize
		if end > len(data) {
			end = len(data)
		}
		copy(fs.ClusterData[currentCluster], data[start:end])
	}
	fs.FatTable[currentCluster] = FAT_EOF

	// Add to directory
	fs.Directory[filename] = DirectoryEntry{
		Name:         filename,
		Size:         len(data),
		FirstCluster: firstCluster,
	}

	fmt.Println("OK")
	return nil
}

// Read a file from the pseudo-FAT (cat s1)
func (fs *FileSystem) cat(filename string) error {

	entry, exists := fs.Directory[filename]
	if !exists {
		return fmt.Errorf("FILE NOT FOUND")
	}

	currentCluster := entry.FirstCluster
	for currentCluster != FAT_EOF {
		fmt.Printf("%s", fs.ClusterData[currentCluster])
		currentCluster = fs.FatTable[currentCluster]
	}
	fmt.Println("OK")
	return nil
}

// Remove a file (rm s1)
func (fs *FileSystem) rm(filename string) error {
	entry, exists := fs.Directory[filename]
	if !exists {
		return fmt.Errorf("FILE NOT FOUND")
	}

	currentCluster := entry.FirstCluster
	for currentCluster != FAT_EOF {
		nextCluster := fs.FatTable[currentCluster]
		fs.FatTable[currentCluster] = FAT_FREE // Mark the cluster as free
		currentCluster = nextCluster
	}

	delete(fs.Directory, filename)
	fmt.Println("OK")
	return nil
}

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
