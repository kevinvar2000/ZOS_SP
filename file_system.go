package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func SaveFileSystem(filename string, fs_format FileSystemFormat, fat1, fat2 FAT) error {

	// **Print the file system details to a file**
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
	for _, val := range fat1 {
		WriteToFile(file, int32(val))
	}

	// **Write FAT2 table at fat2_start position**
	fmt.Println("Seeking to FAT2 start position:", fs_format.fat2_start)
	_, err = file.Seek(int64(fs_format.fat2_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT2 start: %w", err)
	}
	for _, val := range fat2 {
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

func LoadFileSystem(filename string) (FAT, FAT) {

	fmt.Printf("\nLoading file system from '%s'...\n", filename)

	// **Load the file system format from the file**
	fs_format := LoadFormat(filename)
	if fs_format.file_size == 0 {
		fmt.Println("Format file is empty or could not be read.")
		return nil, nil
	}

	// **Open the file for reading**
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil
	}
	defer file.Close()

	fat1 := make(FAT, fs_format.cluster_count)
	fat2 := make(FAT, fs_format.cluster_count)

	// **Read the FAT1 table from the file**
	_, err = file.Seek(int64(fs_format.fat1_start), 0) // Seek to FAT1 start
	if err != nil {
		fmt.Println("Error seeking to FAT1 start:", err)
		return nil, nil
	}
	for i := range fat1 {
		var val int32
		ReadFromFile(file, &val)
		fat1[i] = int(val)
	}

	// **Read the FAT2 table from the file**
	_, err = file.Seek(int64(fs_format.fat2_start), 0) // Seek to FAT2 start
	if err != nil {
		fmt.Println("Error seeking to FAT2 start:", err)
		return nil, nil
	}
	for i := range fat2 {
		var val int32
		ReadFromFile(file, &val)
		fat2[i] = int(val)
	}

	fmt.Println("File system loaded successfully!")

	return fat1, fat2
}

func PrintFileSystem(fat1, fat2 FAT) {

	// **Print the FAT tables**
	fmt.Println("\nFAT1 Table:")
	PrintFAT(fat1)
	fmt.Println("\nFAT2 Table:")
	PrintFAT(fat2)
}

func PrintFAT(fat1 FAT) {
	for i, val := range fat1 {
		fmt.Printf("%d: %d\n", i, val)
	}
}

func SaveFormat(filename string, fs_format FileSystemFormat) {
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

func LoadFormat(filename string) FileSystemFormat {
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

func PrintFormat(fs_format FileSystemFormat) {
	fmt.Printf("\nFile size: %d bytes\n", fs_format.file_size)
	fmt.Printf("FAT size: %d bytes\n", fs_format.fat_size)
	fmt.Printf("FAT cluster count: %d\n", fs_format.fat_cluster_count)
	fmt.Printf("Cluster count: %d\n", fs_format.cluster_count)
	fmt.Printf("FAT1 start: %d\n", fs_format.fat1_start)
	fmt.Printf("FAT2 start: %d\n", fs_format.fat2_start)
	fmt.Printf("Data start: %d\n", fs_format.data_start)
}

func Format(filename string) {

	// **Prompt the user to enter the desired file size**
	var file_size_mb int
	fmt.Print("Enter the desired file size in MB: ")
	fmt.Scanln(&file_size_mb)

	file_size_bytes := file_size_mb * 1024 * 1024

	// **Calculate the file system format**
	fs_format := CalculateFS(file_size_bytes)

	// **Save the file system format to the file**
	SaveFormat(filename, fs_format)

	fat1 := make(FAT, fs_format.cluster_count)
	fat2 := make(FAT, fs_format.cluster_count)

	// **Initialize the FAT table**
	for i := range fat1 {
		fat1[i] = FAT_FREE
		fat2[i] = FAT_FREE
	}

	// **Set the first two entries in the FAT table**
	fat1[0] = FAT_EOF
	fat2[0] = FAT_EOF

	// **Set the entries for the FAT clusters**
	for i := int32(1); i < fs_format.fat_cluster_count; i++ {
		fat1[i] = FAT_FREE
		fat2[i] = FAT_FREE
	}

	// **Save the file system to the file**
	err := SaveFileSystem(filename, fs_format, fat1, fat2)
	if err != nil {
		fmt.Println("Error saving file system:", err)
		return
	}

	// **Find a free cluster for the root directory**
	free_cluster, err := FindFreeCluster(filename)
	if err != nil {
		fmt.Println("Error finding free cluster:", err)
		return
	}

	// **Create the root directory**
	CreateRootDirectory(filename, free_cluster)

	// **Print the file system details**
	PrintFileSystem(fat1, fat2)

	fmt.Printf("File system formatted and saved successfully!\n\n")
}

func CalculateFS(file_size int) FileSystemFormat {

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

func WriteDirectoryEntry(filename string, cluster int32, item DirectoryEntry) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Calculate offset for the cluster and move to position
	offset := int64(cluster * MAX_FILE_NAME * int32(binary.Size(item)))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to cluster: %v", err)
	}

	// Write the directory item
	err = binary.Write(file, binary.LittleEndian, item)
	if err != nil {
		return fmt.Errorf("error writing directory item: %v", err)
	}

	return nil
}

func FindFreeCluster(filename string) (int32, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return -1, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Move to FAT1 start position and read clusters
	fatStart := int64(4) // Assuming FAT starts after a header or root directory
	file.Seek(fatStart, 0)

	var cluster int32
	for i := int32(0); ; i++ {
		err = binary.Read(file, binary.LittleEndian, &cluster)
		if err != nil {
			return -1, fmt.Errorf("error reading FAT: %v", err)
		}
		if cluster == 0 { // 0 indicates a free cluster
			return i, nil
		}
	}
}

func CreateRootDirectory(filename string, free_cluster int32) {

	// **Create the root directory entry**
	rootDir := DirectoryEntry{
		name:          "/",
		size:          0,
		first_cluster: int(free_cluster),
		is_directory:  true,
	}

	// **Write the root directory entry to the file**
	err := WriteDirectoryEntry(filename, free_cluster, rootDir)
	if err != nil {
		fmt.Println("Error writing root directory:", err)
		return
	}

	// **Update the FAT entry for the root directory**
	err = UpdateFatEntry(filename, free_cluster, FAT_EOF)
	if err != nil {
		fmt.Println("Error updating FAT entry for root directory:", err)
		return
	}

	// **Set the current and parent directory for the root directory**
	SetCurrentAndParentDirectory(filename, free_cluster, free_cluster)

	fmt.Println("Root directory set successfully!")

}

func CreateDirectory(filename, dir_name string, parent_cluster int32) {

	// **Check if the directory name is valid**
	if dir_name == "." || dir_name == ".." || len(dir_name) > MAX_FILE_NAME {
		return
	}

	// **Check if the directory name is too long**
	if len(dir_name) > MAX_FILE_NAME {
		fmt.Println("Error: Directory name is too long.")
		return
	}

	// **Check if the directory already exists**
	if CheckIfDirectoryExists(filename, parent_cluster, dir_name) {
		fmt.Println("Error: Directory or file with the name", dir_name, "already exists.")
		return
	}

	// **Find a free cluster for the new directory**
	free_cluster, err := FindFreeCluster(filename)
	if err != nil {
		fmt.Println("Error finding free cluster:", err)
		return
	}

	// **Create the new directory entry**
	new_dir := DirectoryEntry{
		name:          dir_name,
		size:          0,
		first_cluster: int(free_cluster),
		is_directory:  true,
	}

	// **Write the directory entry to the file**
	err = WriteDirectoryEntry(filename, parent_cluster, new_dir)
	if err != nil {
		fmt.Println("Error writing directory entry:", err)
		return
	}

	// **Update the FAT entry for the new directory**
	err = UpdateFatEntry(filename, free_cluster, FAT_EOF)
	if err != nil {
		fmt.Println("Error updating FAT entry:", err)
		return
	}

	// **Set the current and parent directory for the new directory**
	SetCurrentAndParentDirectory(filename, free_cluster, parent_cluster)

	fmt.Printf("Directory '%s' created at cluster %d.\n", dir_name, free_cluster)
}

func SetCurrentAndParentDirectory(filename string, current_cluster, parent_cluster int32) {

	current_entry := DirectoryEntry{
		name:          ".",
		size:          0,
		first_cluster: int(current_cluster),
		is_directory:  true,
	}
	err := WriteDirectoryEntry(filename, current_cluster, current_entry)
	if err != nil {
		return
	}

	parent_entry := DirectoryEntry{
		name:          "..",
		size:          0,
		first_cluster: int(parent_cluster),
		is_directory:  true,
	}
	err = WriteDirectoryEntry(filename, current_cluster, parent_entry)
	if err != nil {
		return
	}

}

func CheckIfDirectoryExists(filename string, parent_cluster int32, dirName string) bool {
	// Implement the logic to check if the directory exists
	// For now, return false as a placeholder
	return false
}

func ReadDirectoryEntrys(filename string, cluster int32) ([]DirectoryEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	offset := int64(cluster * MAX_FILE_NAME * int32(binary.Size(DirectoryEntry{})))
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking to cluster: %v", err)
	}

	var items []DirectoryEntry
	for {
		var item DirectoryEntry
		err = binary.Read(file, binary.LittleEndian, &item)
		if err != nil {
			break
		}
		if item.name != "" {
			items = append(items, item)
		}
	}

	return items, nil
}

func UpdateFatEntry(filename string, cluster, value int32) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fatOffset := int64(4 + cluster*4)
	_, err = file.Seek(fatOffset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT entry: %v", err)
	}

	err = binary.Write(file, binary.LittleEndian, value)
	if err != nil {
		return fmt.Errorf("error updating FAT entry: %v", err)
	}

	return nil
}
