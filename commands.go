package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CopyFile(filename, src, dest string, fs_format FileSystemFormat) {
	// fmt.Println("*** Copying file ***")

	// **Locate the source cluster and name**
	src_cluster, src_name, err := ParsePath(filename, src, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing source path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	// **Locate the source file in the current directory**
	src_entry, err := FindEntry(filename, src_name, src_cluster, fs_format)
	if err != nil {
		// fmt.Println("Error checking source file:", err)
		fmt.Println("FILE NOT FOUND")
		return
	}

	// **Check if source is a directory**
	if src_entry.Is_directory == 1 {
		fmt.Println("Source is a directory and cannot be copied:", src_name)
		return
	}

	// **Locate the destination cluster and name**
	entry_dest_cluster, dest_name, err := ParsePath(filename, dest, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing destination path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	// fmt.Println("Destination Cluster:", entry_dest_cluster, "Destination Name:", dest_name)

	// **Check if destination already exists**
	if CheckIfDirectoryExists(filename, entry_dest_cluster, dest_name, fs_format) {
		fmt.Println("Destination already exists:", dest_name)
		return
	}

	// **Read file contents using the helper function**
	file_contents, err := ReadFileContents(filename, src_entry.First_cluster, src_entry.Size, fs_format)
	if err != nil {
		fmt.Println("Error reading source file contents:", err)
		return
	}

	// **Find the first free cluster for the file**
	dest_cluster, err := FindFreeCluster(filename, fs_format.fat1_start)
	if err != nil {
		fmt.Println("Error finding free cluster:", err)
		return
	}

	if dest_cluster == -1 {
		fmt.Println("Error: Not enough free space in the file system.")
		return
	}

	// fmt.Println("First Free Cluster:", dest_cluster)

	// **Write file contents to the destination clusters**
	err = WriteFileContents(filename, dest_cluster, file_contents, fs_format)
	if err != nil {
		fmt.Println("Error writing file contents:", err)
		return
	}

	// **Create a directory entry for the new file**
	dest_entry := DirectoryEntry{
		Size:          src_entry.Size,
		First_cluster: dest_cluster,
		Is_directory:  src_entry.Is_directory,
	}
	copy(dest_entry.Name[:], dest_name)

	// **Write the new directory entry to the directory**
	err = WriteDirectoryEntry(filename, entry_dest_cluster, dest_entry, fs_format)
	if err != nil {
		fmt.Println("Error writing directory entry:", err)
		return
	}

	// fmt.Println("*** Copy completed successfully ***")
	fmt.Println("OK")
}

func MoveFile(filename, src, dest string, fs_format FileSystemFormat) {
	// fmt.Println("*** Moving file ***")

	src_cluster, src_name, err := ParsePath(filename, src, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing source path:", err)
		fmt.Println("PATH NOT FOUND")
	}

	// fmt.Println("Source Cluster:", src_cluster, "Source Name:", src_name)

	// **Locate the source file in the current directory**
	src_entry, err := FindEntry(filename, src_name, src_cluster, fs_format)
	if err != nil {
		// fmt.Println("Error checking source file:", err)
		fmt.Println("FILE NOT FOUND")
		return
	}

	// **Check if source is a directory**
	if src_entry.Is_directory == 1 {
		fmt.Println("Source is a directory and cannot be moved:", src_name)
		return
	}

	// **Locate the destination cluster and name**
	dest_cluster, dest_name, err := ParsePath(filename, dest, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing destination path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	// fmt.Println("Dest Cluster:", dest_cluster, "Dest Name:", dest_name)

	// **Check if destination already exists**
	var rename bool
	dest_entry, err := FindEntry(filename, dest_name, dest_cluster, fs_format)
	if err != nil {
		fmt.Println("Rename true")
		rename = true
	}

	// **Check if destination already exists**
	if !rename {
		if dest_entry.Is_directory != 1 {
			fmt.Println("Destination is not a directory:", dest_name)
			return
		}
	}

	// **Read the source file contents**
	file_contents, err := ReadFileContents(filename, src_entry.First_cluster, src_entry.Size, fs_format)
	if err != nil {
		fmt.Println("Error reading source file contents:", err)
		return
	}

	// **Find the first free cluster for the file**
	free_cluster, err := FindFreeCluster(filename, fs_format.fat1_start)
	if err != nil {
		fmt.Println("Error finding free cluster:", err)
		return
	}

	if free_cluster == -1 {
		fmt.Println("Error: Not enough free space in the file system.")
		return
	}

	// **Write file contents to the destination clusters**
	err = WriteFileContents(filename, free_cluster, file_contents, fs_format)
	if err != nil {
		fmt.Println("Error writing file contents:", err)
		return
	}
	// **Create a new directory entry for the copied file**
	new_entry := DirectoryEntry{
		Size:          src_entry.Size,
		First_cluster: free_cluster,
		Is_directory:  src_entry.Is_directory,
	}

	if rename {
		copy(new_entry.Name[:], dest_name)
		dest_entry.First_cluster = dest_cluster
	} else {
		copy(new_entry.Name[:], src)
	}

	// **Write the new directory entry to the destination directory**
	// fmt.Println("Writing to cluster:", dest_entry.First_cluster, "Entry:", src)
	err = WriteDirectoryEntry(filename, dest_entry.First_cluster, new_entry, fs_format)
	if err != nil {
		fmt.Println("Error writing directory entry:", err)
		return
	}

	// **Remove the source file from the directory**
	// fmt.Println("Removing in cluster:", current_cluster, "Entry:", src)
	err = RemoveDirectoryEntry(filename, src_cluster, src_name, fs_format)
	if err != nil {
		fmt.Println("Error removing source file:", err)
		return
	}

	// fmt.Println("*** Move completed successfully ***")
	fmt.Println("OK")
}

func RemoveFile(filename, file string, fs_format FileSystemFormat) {

	// fmt.Println("*** Removing file ***")

	file_cluster, file_name, err := ParsePath(filename, file, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing file path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	// **Locate the source file in the current directory**
	entry, err := FindEntry(filename, file_name, file_cluster, fs_format)
	if err != nil {
		// fmt.Println("Error checking file:", err)
		fmt.Println("FILE NOT FOUND")
		return
	}

	// **Check if file exists**
	if IsZeroEntry(entry) {
		fmt.Println("File not found:", file_name)
		return
	}

	// **Remove the file**
	err = RemoveDirectoryEntry(filename, file_cluster, file_name, fs_format)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	// fmt.Println("File removed successfully:", file)
	fmt.Println("OK")
}

func MakeDirectory(dir_name, filename string, fs_format FileSystemFormat) {

	CreateDirectory(filename, dir_name, fs_format)
	fmt.Println("OK")
}

func RemoveDirectory(filename, dir_name string, fs_format FileSystemFormat) {

	dir_cluster, dir_name, err := ParsePath(filename, dir_name, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing directory path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	err = RemoveDirectoryEntry(filename, dir_cluster, dir_name, fs_format)
	if err != nil {
		fmt.Println("Error removing directory:", err)
		return
	}
	fmt.Println("OK")
}

func PrintDirectoryContents(filename, src string, fs_format FileSystemFormat) {

	current_cluster := GetCurrentCluster()

	// fmt.Println("Current Cluster:", current_cluster)

	if src != "" {

		// fmt.Println("Changing directory to:", src)

		var err error
		current_cluster, _, err = ParsePath(filename, src, fs_format, false)
		if err != nil {
			// fmt.Println("Error parsing path:", err)
			fmt.Println("PATH NOT FOUND")
			return
		}
		// fmt.Println("New Cluster:", current_cluster)
	}

	// fmt.Println("Current Cluster before read:", current_cluster)

	dir_entries, err := ReadDirectoryEntries(filename, current_cluster, fs_format)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("Directory Contents:")
	fmt.Printf("%-20s %-10s %-15s %-15s\n", "Name", "Size", "First Cluster", "Is Directory")

	for _, entry := range dir_entries {

		if !IsZeroEntry(entry) {

			dir_name_str := string(bytes.Trim(entry.Name[:], "\x00"))
			fmt.Printf("%-20s %-10d %-15d %-15d\n", dir_name_str, entry.Size, entry.First_cluster, entry.Is_directory)
		}
	}
}

func PrintFileContents(filename, file string, fs_format FileSystemFormat) {

	file_cluster, file_name, err := ParsePath(filename, file, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing file path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	// **Locate the file in the current directory**
	entry, err := FindEntry(filename, file_name, file_cluster, fs_format)
	if err != nil {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// **Read the file contents**
	file_contents, err := ReadFileContents(filename, entry.First_cluster, entry.Size, fs_format)
	if err != nil {
		fmt.Println("Error reading file contents:", err)
		return
	}

	// **Print the file contents**
	fmt.Println(string(file_contents))

	// fmt.Println("File Contents:")
	// fmt.Println(string(file_contents))
}

func ChangePath(filename, path string, fs_format FileSystemFormat) {

	// fmt.Println("Changing directory to:", path)

	// Check if the path is absolute or relative
	var start_cluster int32
	if path == "/" {
		// fmt.Println("Root directory")
		// Root directory case
		SetCurrentCluster(fs_format.data_start / CLUSTER_SIZE)
		SetCurrentPath(path)
		fmt.Println("OK")
		return
	} else if path[0] == '/' {
		// fmt.Println("Absolute path")
		// Absolute path, start from root cluster
		start_cluster = GetCurrentCluster()
		path = path[1:] // Remove leading '/'
	} else {
		// fmt.Println("Relative path")
		// Relative path, start from the current directory cluster
		start_cluster = GetCurrentCluster()
	}

	// fmt.Println("Start Cluster:", start_cluster)

	// Split the path into directory components
	path_components := strings.Split(path, "/")
	current_cluster := start_cluster

	// Traverse each component in the path
	for _, component := range path_components {

		if component == "." || component == "" {
			// Stay in the current directory
			// fmt.Println("Staying in the current directory")
			continue
		}

		if component == ".." {
			// Move up to the parent directory
			// fmt.Println("Moving up to the parent directory")
			current_cluster = GetParentCluster(filename, current_cluster, fs_format)
			continue
		}

		// Search for the component in the current directory entries
		dir_entries, err := ReadDirectoryEntries(filename, current_cluster, fs_format)
		if err != nil {
			fmt.Println("Error reading directory entries:", err)
			return
		}

		found := false
		for _, entry := range dir_entries {
			entryName := bytes.Trim(entry.Name[:], "\x00")
			if entry.Is_directory == 1 && string(entryName) == component {
				// Found the directory; update the current cluster
				current_cluster = entry.First_cluster
				found = true
				// fmt.Println("Found directory:", component)
				// fmt.Println("New Cluster:", current_cluster)
				// fmt.Println()
				break
			}
		}

		if !found {
			fmt.Println("PATH NOT FOUND")
			return
		}
	}

	// Update the current directory cluster if path traversal succeeded
	SetCurrentCluster(current_cluster)
	SetCurrentPath(path)
	// fmt.Println("Directory changed to:", path)
	fmt.Println("OK")
}

func PrintCurrentPath() {

	GetCurrentPath()
}

func PrintInformation(filename, src string, fs_format FileSystemFormat) {

	src_cluster, src_name, err := ParsePath(filename, src, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing file path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	src_entry, err := FindEntry(filename, src_name, src_cluster, fs_format)
	if err != nil {
		fmt.Println("Error checking file:", err)
		return
	}

	// fmt.Println("Source cluster:", src_cluster, "Source name:", src_name)
	// fmt.Println("Source entry cluster:", src_entry.First_cluster)

	fmt.Print(src_name, ": ")
	current_cluster := src_entry.First_cluster
	for current_cluster != FAT_EOF {

		fmt.Printf("%d ", current_cluster)

		next_cluster, err := ReadFatEntry(filename, current_cluster, fs_format)
		if err != nil {
			fmt.Println("Error reading FAT entry:", err)
			return
		}

		if next_cluster == FAT_EOF {
			break
		}

		current_cluster = next_cluster

	}
	fmt.Println()
}

func Incp(filename string, src string, dest string, fs_format FileSystemFormat) {

	// **Open the source file for reading**
	file, err := os.Open(src)
	if err != nil {
		fmt.Println("FILE NOT FOUND")
		// fmt.Println("Error opening source file:", err)
		return
	}
	defer file.Close()

	// **Get the size of the source file**
	file_info, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return
	}
	file_size := int32(file_info.Size())

	// fmt.Println("Source File Size:", file_size)

	// **Read the source file's contents**
	file_contents := make([]byte, file_size)
	_, err = file.Read(file_contents)
	if err != nil {
		fmt.Println("Error reading source file:", err)
		return
	}

	// fmt.Println("File contents size:", len(file_contents))
	// fmt.Println("File contents:", string(file_contents))
	// fmt.Println("File contents:", file_contents)

	// **Parse the destination path**
	dest_cluster, dest_name, err := ParsePath(filename, dest, fs_format, true)
	if err != nil {
		fmt.Println("PATH NOT FOUND")
		// fmt.Println("Error parsing destination path:", err)
		return
	}

	// fmt.Println("Destination Cluster:", dest_cluster, "Destination Name:", dest_name)

	// **Check if a file with the same name already exists**
	dir_entries, err := ReadDirectoryEntries(filename, dest_cluster, fs_format)
	if err != nil {
		fmt.Println("Error reading directory entries:", err)
		return
	}

	for _, entry := range dir_entries {
		if string(bytes.Trim(entry.Name[:], "\x00")) == dest_name {
			fmt.Println("Destination file already exists:", dest_name)
			return
		}
	}

	// **Find the first free cluster for the file**
	first_cluster, err := FindFreeCluster(filename, fs_format.fat1_start)
	if err != nil {
		fmt.Println("Error finding free cluster:", err)
		return
	}

	if first_cluster == -1 {
		fmt.Println("Error: Not enough free space in the file system.")
		return
	}

	// fmt.Println("First Free Cluster:", first_cluster)

	// **Write the file data into the VFS**
	err = WriteFileContents(filename, first_cluster, file_contents, fs_format)
	if err != nil {
		fmt.Println("Error writing file contents:", err)
		return
	}

	// **Create a directory entry for the new file**
	new_entry := DirectoryEntry{
		Size:          file_size,
		First_cluster: first_cluster,
		Is_directory:  0, // 0 indicates a file
	}
	copy(new_entry.Name[:], dest_name)

	// **Write the new directory entry to the directory**
	err = WriteDirectoryEntry(filename, dest_cluster, new_entry, fs_format)
	if err != nil {
		fmt.Println("Error writing directory entry:", err)
		return
	}

	// fmt.Println("File successfully copied into the virtual file system as:", dest)
	fmt.Println("OK")
}

func Outcp(filename string, src string, dest string, fs_format FileSystemFormat) {

	// **Parse the source path and locate the source file**
	src_cluster, src_name, err := ParsePath(filename, src, fs_format, true)
	if err != nil {
		fmt.Println("Error parsing source path:", err)
		return
	}

	// fmt.Println("Source Cluster:", src_cluster)

	// **Locate the source file in the directory**
	dir_entries, err := ReadDirectoryEntries(filename, src_cluster, fs_format)
	if err != nil {
		fmt.Println("Error reading directory entries:", err)
		return
	}

	var src_entry DirectoryEntry
	found := false
	for _, entry := range dir_entries {
		if string(bytes.Trim(entry.Name[:], "\x00")) == src_name {
			src_entry = entry
			found = true
			break
		}
	}

	if !found {
		fmt.Println("FILE NOT FOUND")
		// fmt.Println("Source file not found:", src)
		return
	}

	// **Check if source is a directory**
	if src_entry.Is_directory == 1 {
		fmt.Println("Source is a directory and cannot be copied:", src)
		return
	}

	// **Read the file contents from the VFS using ReadFileContents**
	file_contents, err := ReadFileContents(filename, src_entry.First_cluster, src_entry.Size, fs_format)
	if err != nil {
		fmt.Println("Error reading source file from VFS:", err)
		return
	}

	// **Open the destination file for writing**
	destFile, err := os.Create(dest)
	if err != nil {
		fmt.Println("PATH NOT FOUND")
		// fmt.Println("Error creating destination file:", err)
		return
	}
	defer destFile.Close()

	// **Write the file contents to the destination file**
	_, err = destFile.Write(file_contents)
	if err != nil {
		fmt.Println("Error writing to destination file:", err)
		return
	}

	// fmt.Println("File successfully copied from VFS to external destination:", dest)
	fmt.Println("OK")
}

func LoadFile(filename, script string, fs_format FileSystemFormat) {

	// **Read the commands from the script file**
	data, err := os.ReadFile(script)
	if err != nil {
		fmt.Println("FILE NOT FOUND")
		return
	}

	// **Split the commands by newline**
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {

		// **Print empty lines**
		if line == "" || line == "\n" || line == "\r" || line == "\r\n" || strings.TrimSpace(line) == "" {
			fmt.Println()
			continue
		}

		// **Print the comment**
		if strings.HasPrefix(line, "#") {
			fmt.Println(line)
			continue
		}

		// **Split the command by space**
		words := strings.Fields(line)

		var command, arg1, arg2 string
		if len(words) == 1 {
			command = words[0]
		} else if len(words) == 2 {
			command = words[0]
			arg1 = words[1]
		} else if len(words) == 3 {
			command = words[0]
			arg1 = words[1]
			arg2 = words[2]
		}

		fmt.Println("Executing:", command, arg1, arg2)
		ExecuteCommand(filename, command, arg1, arg2, fs_format)
	}

	// fmt.Println("OK")
}

func FormatFileCmd(filename string, size int) {
	Format(filename, size)
	fmt.Println("OK")
}

func BugTest(filename, bug_file string, fs_format FileSystemFormat) {

	bug_file_cluster, bug_file_name, err := ParsePath(filename, bug_file, fs_format, true)
	if err != nil {
		fmt.Println("Error parsing file path:", err)
		return
	}

	// Read all directory entries in the current directory
	dir_entries, err := ReadDirectoryEntries(filename, bug_file_cluster, fs_format)
	if err != nil {
		fmt.Println("Error reading directory entries:", err)
		return
	}

	// Search for the specified file in the directory
	var start_cluster int32 = -1
	for _, entry := range dir_entries {
		if !IsZeroEntry(entry) { // Skip empty entries
			entryName := string(bytes.Trim(entry.Name[:], "\x00"))
			if entryName == bug_file_name {
				start_cluster = entry.First_cluster
				break
			}
		}
	}

	// If the file is not found, display an error message
	if start_cluster == -1 {
		fmt.Printf("Error: File '%s' not found or is not a file.\n", bug_file)
		return
	}

	// Update the FAT entry to mark the file as corrupted (FAT_BAD_CLUSTER)
	err = UpdateFatEntry(filename, start_cluster, FAT_BAD, fs_format)
	if err != nil {
		fmt.Printf("Error marking file '%s' as corrupted: %v\n", bug_file, err)
		return
	}

	fmt.Printf("Marked file '%s' as corrupted (FAT_BAD_CLUSTER).\n", bug_file)
}

func CheckForBugs(filename string, fs_format FileSystemFormat) {

	// **Load the FAT tables from the file system**
	fat1, fat2 := LoadFileSystem(filename)

	// **Check for bad clusters in the FAT tables**
	for i := 0; i < len(fat1); i++ {
		if fat1[i] == FAT_BAD || fat2[i] == FAT_BAD {
			fmt.Printf("Bad cluster found in FAT table: Cluster %d\n", i)
		}
	}

	fmt.Println("OK")
}

func PrintHelp() {
	fmt.Println("Commands:")
	fmt.Println("cp - Copy the file")
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
	fmt.Println("bug - Bug test")
	fmt.Println("check - Check for bugs")
	fmt.Println("print - Print the FAT tables to the file")
	fmt.Println("help - Print the help")
	fmt.Println("exit - Exit the program")
	fmt.Println()
}

func ExecuteCommand(filename, command, arg1, arg2 string, fs_format FileSystemFormat) {

	switch command {
	case "cp":
		if arg1 == "" || arg2 == "" {
			fmt.Println("Source and destination paths are required for copy.")
			return
		}
		CopyFile(filename, arg1, arg2, fs_format)
	case "mv":
		if arg1 == "" || arg2 == "" {
			fmt.Println("Source and destination paths are required for move.")
			return
		}
		MoveFile(filename, arg1, arg2, fs_format)
	case "rm":
		if arg1 == "" {
			fmt.Println("File path is required for remove.")
			return
		}
		RemoveFile(filename, arg1, fs_format)
	case "mkdir":
		if arg1 == "" {
			fmt.Println("Directory name is required for mkdir.")
			return
		}
		MakeDirectory(arg1, filename, fs_format)
	case "rmdir":
		if arg1 == "" {
			fmt.Println("Directory name is required for rmdir.")
			return
		}
		RemoveDirectory(filename, arg1, fs_format)
	case "ls":
		PrintDirectoryContents(filename, arg1, fs_format)
	case "cat":
		if arg1 == "" {
			fmt.Println("File path is required for cat.")
			return
		}
		PrintFileContents(filename, arg1, fs_format)
	case "cd":
		if arg1 == "" {
			fmt.Println("Path is required for cd.")
			return
		}
		ChangePath(filename, arg1, fs_format)
	case "pwd":
		PrintCurrentPath()
	case "info":
		if arg1 == "" {
			fmt.Println("File path is required for info.")
			return
		}
		PrintInformation(filename, arg1, fs_format)
	case "incp":
		if arg1 == "" || arg2 == "" {
			fmt.Println("Source and destination paths are required for incp.")
			return
		}
		Incp(filename, arg1, arg2, fs_format)
	case "outcp":
		if arg1 == "" || arg2 == "" {
			fmt.Println("Source and destination paths are required for outcp.")
			return
		}
		Outcp(filename, arg1, arg2, fs_format)
	case "load":
		if arg1 == "" {
			fmt.Println("Script file path is required for load.")
			return
		}
		LoadFile(filename, arg1, fs_format)
	case "format":
		if arg1 == "" {
			fmt.Println("Size is required for format.")
			return
		}
		size, err := strconv.Atoi(arg1)
		if err != nil {
			fmt.Println("Invalid size:", arg1)
			return
		}
		FormatFileCmd(filename, size)
	case "bug":
		if arg1 == "" {
			fmt.Println("File name is required for bug.")
			return
		}
		BugTest(filename, arg1, fs_format)
	case "check":
		CheckForBugs(filename, fs_format)
	case "print":
		fat1, fat2 := LoadFileSystem(filename)
		PrintFileSystem(fat1, fat2, "fats.txt")
	case "help":
		PrintHelp()
	case "exit":
	case "quit":
	case "q":
		fmt.Println("Exiting the file system simulator.")
		return
	default:
		fmt.Println("Invalid command")
	}
}
