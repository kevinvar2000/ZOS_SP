package main

import (
	"fmt"
	"os"
)

func SaveFileSystem(filename string, fs_format FileSystemFormat) error {

	fat_table := getFAT()

	// **Print the file system details to a file**
	// PrintFileSystem(fs, "fs_saving.txt")
	fmt.Printf("\nSaving file system to '%s'...\n", filename)

	// **Open the file for writing**
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// **Write FAT1 table at fat1_start position**
	fmt.Println("Seeking to FAT1 start position:", fs_format.fat1_start)
	_, err = file.Seek(int64(fs_format.fat1_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT1 start: %w", err)
	}
	for _, val := range fat_table {
		WriteToFile(file, int32(val))
	}

	// **Write FAT2 table at fat2_start position**
	fmt.Println("Seeking to FAT2 start position:", fs_format.fat2_start)
	_, err = file.Seek(int64(fs_format.fat2_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT2 start: %w", err)
	}
	for _, val := range fat_table {
		WriteToFile(file, int32(val))
	}

	// **Zero out the data starting at data_start**
	remainingSize := fs_format.file_size - fs_format.data_start
	zeroBuffer := make([]byte, remainingSize)
	fmt.Println("Zero buffer size:", len(zeroBuffer))

	// **Write the zero buffer to the data section**
	fmt.Println("Seeking to data start position:", fs_format.data_start)
	_, err = file.Seek(int64(fs_format.data_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to data start: %w", err)
	}
	_, err = file.Write(zeroBuffer)
	if err != nil {
		return fmt.Errorf("error writing zeros to data section: %w", err)
	}

	fmt.Printf("File system saved successfully!\n\n")
	return nil
}

func LoadFileSystem(filename string) {

	fmt.Printf("\nLoading file system from '%s'...\n", filename)

	// **Load the file system format from the file**
	fs_format := LoadFormatFile(filename)
	if fs_format.file_size == 0 {
		fmt.Println("Format file is empty or could not be read.")
		return
	}

	// **Open the file for reading**
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fat_table := getFAT()

	// **Read the FAT1 table from the file**
	_, err = file.Seek(int64(fs_format.fat1_start), 0) // Seek to FAT1 start
	if err != nil {
		fmt.Println("Error seeking to FAT1 start:", err)
		return
	}
	for i := range fat_table {
		var val int32
		ReadFromFile(file, &val)
		fat_table[i] = int(val)
	}

	// **Read the data clusters from the file**
	_, err = file.Seek(int64(fs_format.data_start), 0)
	if err != nil {
		fmt.Println("Error seeking to data start:", err)
		return
	}

	cluster_data := make([][]byte, fs_format.cluster_count)
	for i := range fs_format.cluster_count {
		var val int32
		ReadFromFile(file, &val)
		cluster_data[i] = make([]byte, CLUSTER_SIZE)
		// cluster_data[i] = val
	}

	fmt.Println("File system loaded successfully!")
	return
}

func SaveFormatFile(filename string, fs_format FileSystemFormat) {
	fmt.Printf("\nSaving format of the file system to '%s'...\n", filename)

	// **Open the file for writing**
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// **Write the file system format to the file**
	WriteToFile(file, fs_format.file_size)
	fmt.Printf("File size written: %d bytes\n", fs_format.file_size)

	// **Write the FAT size**
	WriteToFile(file, fs_format.fat_size)
	fmt.Printf("FAT size written: %d bytes\n", fs_format.fat_size)

	// **Write the number of FAT clusters**
	WriteToFile(file, fs_format.fat_cluster_count)
	fmt.Printf("FAT cluster count written: %d\n", fs_format.fat_cluster_count)

	// **Write the total number of data clusters**
	WriteToFile(file, fs_format.cluster_count)
	fmt.Printf("Cluster count written: %d\n", fs_format.cluster_count)

	// **Write the starting positions**
	WriteToFile(file, fs_format.fat1_start)
	fmt.Printf("FAT1 starts at: %d\n", fs_format.fat1_start)

	WriteToFile(file, fs_format.fat2_start)
	fmt.Printf("FAT2 starts at: %d\n", fs_format.fat2_start)

	WriteToFile(file, fs_format.data_start)
	fmt.Printf("Data starts at: %d\n", fs_format.data_start)

	fmt.Printf("File system format saved successfully!\n\n")
}

func LoadFormatFile(filename string) FileSystemFormat {
	fmt.Printf("Loading format of the file system from '%s'...\n", filename)

	// **Initialize the file system format**
	fs_format := FileSystemFormat{}

	// **Open the file for reading**
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return FileSystemFormat{}
	}
	defer file.Close()

	// **Read the file system format from the file**
	ReadFromFile(file, &fs_format.file_size)
	fmt.Printf("File size read: %d bytes\n", fs_format.file_size)

	// **Read the FAT size**
	ReadFromFile(file, &fs_format.fat_size)
	fmt.Printf("FAT size read: %d bytes\n", fs_format.fat_size)

	// **Read the number of FAT clusters**
	ReadFromFile(file, &fs_format.fat_cluster_count)
	fmt.Printf("FAT cluster count read: %d\n", fs_format.fat_cluster_count)

	// **Read the total number of data clusters**
	ReadFromFile(file, &fs_format.cluster_count)
	fmt.Printf("Cluster count read: %d\n", fs_format.cluster_count)

	// **Read the starting positions**
	ReadFromFile(file, &fs_format.fat1_start)
	fmt.Printf("FAT1 starts at: %d\n", fs_format.fat1_start)

	// **Read the starting positions**
	ReadFromFile(file, &fs_format.fat2_start)
	fmt.Printf("FAT2 starts at: %d\n", fs_format.fat2_start)

	// **Read the starting positions**
	ReadFromFile(file, &fs_format.data_start)
	fmt.Printf("Data starts at: %d\n", fs_format.data_start)

	fmt.Printf("File system format loaded successfully!\n\n")
	return fs_format
}

// Format the file with the desired size
func FSFormatFile(filename string) {

	// **Prompt the user to enter the desired file size**
	var file_size int
	fmt.Print("Enter the desired file size in bytes: ")
	fmt.Scanln(&file_size)

	file_size_mb := file_size * 1024 * 1024 // Convert to MB

	// **Calculate the file system format**
	fs_format := CalculateFSFormat(file_size_mb)

	// **Save the file system format to the file**
	SaveFormatFile(filename, fs_format)

	// **Save the file system to the file**
	err := SaveFileSystem(filename, fs_format)
	if err != nil {
		fmt.Println("Error saving file system:", err)
		return
	}

	CreateDirectory(filename, "/")

	fmt.Printf("File system formatted and saved successfully!\n\n")
}

func CalculateFSFormat(file_size int) FileSystemFormat {

	fmt.Printf("File size: %d bytes\n", file_size)

	// **Calculate the number of clusters based on the file size**
	cluster_count := int(file_size / CLUSTER_SIZE)
	fmt.Printf("Cluster count: %d\n", cluster_count)

	// **Calculate the FAT size and number of FAT clusters**
	fat_size := cluster_count * FAT_ENTRY
	fat_cluster_count := (fat_size + CLUSTER_SIZE - 1) / CLUSTER_SIZE

	fmt.Printf("FAT size: %d bytes\n", fat_size)
	fmt.Printf("FAT cluster count: %d\n", fat_cluster_count)

	// **Calculate the starting positions**
	fat1_start := CLUSTER_SIZE
	fat2_start := fat1_start + fat_cluster_count*CLUSTER_SIZE
	data_start := fat2_start + fat_cluster_count*CLUSTER_SIZE

	fmt.Printf("FAT1 starts at: %d\n", fat1_start)
	fmt.Printf("FAT2 starts at: %d\n", fat2_start)
	fmt.Printf("Data starts at: %d\n", data_start)

	// **Initialize the file system format**
	fs_format := FileSystemFormat{
		file_size:         int32(file_size),
		fat_size:          int32(fat_size),
		fat_cluster_count: int32(fat_cluster_count),
		cluster_count:     int32(cluster_count),
		fat1_start:        int32(fat1_start),
		fat2_start:        int32(fat2_start),
		data_start:        int32(data_start),
	}

	return fs_format
}

func CreateDirectory(filename string, path string) {

	if path == "." || path == ".." {
		fmt.Println("Invalid directory name!")
		return
	}

	// **Check if the root directory is being created**
	if path == "/" {

		// **Create the root directory entry**
		// rootDir := DirectoryEntry{
		// 	name:          "/",
		// 	size:          0,
		// 	first_cluster: 0,
		// 	is_directory:  true,
		// }

		// TODO: find free cluster and write to it

	} else {

		// **Create the directory entry**
		// dir := DirectoryEntry{
		// 	name:          path,
		// 	size:          0,
		// 	first_cluster: 0,
		// 	is_directory:  true,
		// }

		fmt.Printf("Directory '%s' created successfully!\n", path)
	}

}

func getFAT() FAT {

	// **Initialize the FAT table**
	FAT := make(FAT, 10)

	return FAT
}
