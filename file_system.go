package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

var current_cluster int32
var current_path string = "/"

func SaveFileSystem(filename string, fs_format FileSystemFormat, fat1, fat2 FAT) error {

	// **Print the file system details to a file**
	// fmt.Printf("\nSaving file system to '%s'...\n", filename)

	// **Open the file for writing**
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// **Write FAT1 table at fat1_start position**
	// fmt.Println("Seeking to FAT1 start position:", fs_format.fat1_start)
	_, err = file.Seek(int64(fs_format.fat1_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT1 start: %w", err)
	}
	for _, val := range fat1 {
		WriteToFile(file, int32(val))
	}

	// **Write FAT2 table at fat2_start position**
	// fmt.Println("Seeking to FAT2 start position:", fs_format.fat2_start)
	_, err = file.Seek(int64(fs_format.fat2_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT2 start: %w", err)
	}
	for _, val := range fat2 {
		WriteToFile(file, int32(val))
	}

	// **Zero out the data starting at data_start**
	remaining_size := fs_format.file_size - fs_format.data_start
	zero_buffer := make([]byte, remaining_size)
	// fmt.Println("Zero buffer size:", len(zero_buffer))

	// **Write the zero buffer to the data section**
	// fmt.Println("Seeking to data start position:", fs_format.data_start)
	_, err = file.Seek(int64(fs_format.data_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to data start: %w", err)
	}
	_, err = file.Write(zero_buffer)
	if err != nil {
		return fmt.Errorf("error writing zeros to data section: %w", err)
	}

	// fmt.Printf("File system saved successfully!\n\n")
	return nil
}

func LoadFileSystem(filename string) (FAT, FAT) {

	// fmt.Printf("\nLoading file system from '%s'...\n", filename)

	// **Load the file system format from the file**
	fs_format := LoadFormat(filename)
	if fs_format.file_size == 0 {
		// fmt.Println("Format file is empty or could not be read.")
		return nil, nil
	}

	// **Open the file for reading**
	file, err := os.Open(filename)
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return nil, nil
	}
	defer file.Close()

	fat1 := make(FAT, fs_format.cluster_count)
	fat2 := make(FAT, fs_format.cluster_count)

	// **Read the FAT1 table from the file**
	_, err = file.Seek(int64(fs_format.fat1_start), 0) // Seek to FAT1 start
	if err != nil {
		// fmt.Println("Error seeking to FAT1 start:", err)
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
		// fmt.Println("Error seeking to FAT2 start:", err)
		return nil, nil
	}
	for i := range fat2 {
		var val int32
		ReadFromFile(file, &val)
		fat2[i] = int(val)
	}

	// fmt.Println("File system loaded successfully!")

	return fat1, fat2
}

func PrintFileSystem(fat1, fat2 FAT, filename string) error {

	// **Open the file for writing**
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// **Print the FAT1 table to the file**
	_, err = file.WriteString("\nFAT1 Table:\n")
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	for i, val := range fat1 {
		_, err = file.WriteString(fmt.Sprintf("%d: %d\n", i, val))
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	// **Print the FAT2 table to the file**
	_, err = file.WriteString("\nFAT2 Table:\n")
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	for i, val := range fat2 {
		_, err = file.WriteString(fmt.Sprintf("%d: %d\n", i, val))
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	// fmt.Println("File system details printed to file successfully!")
	return nil
}

func SaveFormat(filename string, fs_format FileSystemFormat) {
	// fmt.Printf("\nSaving format of the file system to '%s'...\n", filename)

	// **Open the file for writing**
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// **Write the file system format to the file**
	WriteToFile(file, fs_format.file_size)
	// fmt.Printf("File size written: %d bytes\n", fs_format.file_size)

	// **Write the FAT size**
	WriteToFile(file, fs_format.fat_size)
	// fmt.Printf("FAT size written: %d bytes\n", fs_format.fat_size)

	// **Write the number of FAT clusters**
	WriteToFile(file, fs_format.fat_cluster_count)
	// fmt.Printf("FAT cluster count written: %d\n", fs_format.fat_cluster_count)

	// **Write the total number of data clusters**
	WriteToFile(file, fs_format.cluster_count)
	// fmt.Printf("Cluster count written: %d\n", fs_format.cluster_count)

	// **Write the starting positions**
	WriteToFile(file, fs_format.fat1_start)
	// fmt.Printf("FAT1 starts at: %d\n", fs_format.fat1_start)

	WriteToFile(file, fs_format.fat2_start)
	// fmt.Printf("FAT2 starts at: %d\n", fs_format.fat2_start)

	WriteToFile(file, fs_format.data_start)
	// fmt.Printf("Data starts at: %d\n", fs_format.data_start)

	// fmt.Printf("File system format saved successfully!\n\n")
}

func LoadFormat(filename string) FileSystemFormat {
	// fmt.Printf("Loading format of the file system from '%s'...\n", filename)

	// **Initialize the file system format**
	fs_format := FileSystemFormat{}

	// **Open the file for reading**
	file, err := os.Open(filename)
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return FileSystemFormat{}
	}
	defer file.Close()

	// **Read the file system format from the file**
	ReadFromFile(file, &fs_format.file_size)
	// fmt.Printf("File size read: %d bytes\n", fs_format.file_size)

	// **Read the FAT size**
	ReadFromFile(file, &fs_format.fat_size)
	// fmt.Printf("FAT size read: %d bytes\n", fs_format.fat_size)

	// **Read the number of FAT clusters**
	ReadFromFile(file, &fs_format.fat_cluster_count)
	// fmt.Printf("FAT cluster count read: %d\n", fs_format.fat_cluster_count)

	// **Read the total number of data clusters**
	ReadFromFile(file, &fs_format.cluster_count)
	// fmt.Printf("Cluster count read: %d\n", fs_format.cluster_count)

	// **Read the starting positions**
	ReadFromFile(file, &fs_format.fat1_start)
	// fmt.Printf("FAT1 starts at: %d\n", fs_format.fat1_start)

	// **Read the starting positions**
	ReadFromFile(file, &fs_format.fat2_start)
	// fmt.Printf("FAT2 starts at: %d\n", fs_format.fat2_start)

	// **Read the starting positions**
	ReadFromFile(file, &fs_format.data_start)
	// fmt.Printf("Data starts at: %d\n", fs_format.data_start)

	// fmt.Printf("File system format loaded successfully!\n\n")
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

func Format(filename string, file_size_mb int) {

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
	for i := int32(1); i < 2*fs_format.fat_cluster_count+1; i++ {
		fat1[i] = FAT_EOF
		fat2[i] = FAT_EOF
	}

	// **Save the file system to the file**
	err := SaveFileSystem(filename, fs_format, fat1, fat2)
	if err != nil {
		// fmt.Println("Error saving file system:", err)
		return
	}

	// **Find a free cluster for the root directory**
	free_cluster, err := FindFreeCluster(filename, fs_format.fat1_start)
	if err != nil {
		// fmt.Println("Error finding free cluster:", err)
		return
	}

	// **Create the root directory**
	CreateRootDirectory(filename, free_cluster, fs_format)

	// fmt.Printf("File system formatted and saved successfully!\n\n")
}

func CalculateFS(file_size int) FileSystemFormat {

	// fmt.Printf("File size: %d bytes\n", file_size)

	// **Calculate the number of clusters based on the file size**
	cluster_count := int(file_size / CLUSTER_SIZE)
	// fmt.Printf("Cluster count: %d\n", cluster_count)

	// **Calculate the FAT size and number of FAT clusters**
	fat_size := cluster_count * FAT_ENTRY
	fat_cluster_count := (fat_size + CLUSTER_SIZE - 1) / CLUSTER_SIZE

	// fmt.Printf("FAT size: %d bytes\n", fat_size)
	// fmt.Printf("FAT cluster count: %d\n", fat_cluster_count)

	// **Calculate the starting positions**
	fat1_start := CLUSTER_SIZE
	fat2_start := fat1_start + fat_cluster_count*CLUSTER_SIZE
	data_start := fat2_start + fat_cluster_count*CLUSTER_SIZE

	// fmt.Printf("FAT1 starts at: %d\n", fat1_start)
	// fmt.Printf("FAT2 starts at: %d\n", fat2_start)
	// fmt.Printf("Data starts at: %d\n", data_start)

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

func WriteDirectoryEntry(filename string, cluster int32, dir_entry DirectoryEntry, fs_format FileSystemFormat) error {

	// fmt.Println("*** Writing directory entry ***")

	// **Read the directory entries from the cluster**
	dir_entries, err := ReadDirectoryEntries(filename, cluster, fs_format)
	if err != nil {
		return fmt.Errorf("error reading directory entries: %v", err)
	}

	// **Find the first empty slot in the directory**
	empty_index := -1
	for i, entry := range dir_entries {
		if IsZeroEntry(entry) {
			empty_index = i
			break
		}
	}

	// **Check if an empty slot was found**
	if empty_index == -1 {
		return fmt.Errorf("no empty directory slot available in cluster %d", cluster)
	}

	// **Write the directory entry to the empty slot**
	dir_entries[empty_index] = dir_entry

	// **Calculate the data cluster position for the directory entry**
	data_cluster := cluster - 2*fs_format.fat_cluster_count - 1
	cluster_offset := int64(fs_format.data_start + data_cluster*CLUSTER_SIZE)

	// **Write the directory entries back to the file**
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	_, err = file.Seek(cluster_offset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to cluster offset: %v", err)
	}

	for _, entry := range dir_entries {
		err = binary.Write(file, binary.LittleEndian, entry)
		if err != nil {
			return fmt.Errorf("error writing directory entry: %v", err)
		}
	}

	// fmt.Println("Directory entry written successfully!")
	// fmt.Println()
	return nil
}

func FindFreeCluster(filename string, fat_start int32) (int32, error) {

	// fmt.Println("*** Finding free cluster ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return -1, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// **Seek to the start of the FAT table**
	_, err = file.Seek(int64(fat_start), 0)
	if err != nil {
		return -1, fmt.Errorf("error seeking to FAT start: %v", err)
	}

	// **Find the first free cluster in the FAT table**
	var cluster int32
	for i := int32(0); ; i++ {

		err = binary.Read(file, binary.LittleEndian, &cluster)
		if err != nil {
			return -1, fmt.Errorf("error reading FAT: %v", err)
		}

		if cluster == -1 {
			// fmt.Println("Free cluster found at:", i)
			return i, nil
		}
	}
}

func CreateRootDirectory(filename string, free_cluster int32, fs_format FileSystemFormat) {

	// fmt.Println("*** Creating root directory ***")

	// **Create the root directory entry**
	// rootDir := DirectoryEntry{
	// 	Name:          [MAX_FILE_NAME]byte{'/'},
	// 	Size:          0,
	// 	First_cluster: free_cluster,
	// 	Is_directory:  1,
	// }

	// fmt.Println("Root directory:", rootDir)
	// fmt.Println("Root directory size:", binary.Size(rootDir))
	// fmt.Println("Free cluster:", free_cluster)

	// **Write the root directory entry to the file**
	// err := WriteDirectoryEntry(filename, free_cluster, rootDir, fs_format)
	// if err != nil {
	// 	// fmt.Println("Error writing root directory:", err)
	// 	return
	// }

	// **Update the FAT entry for the root directory**
	err := UpdateFatEntry(filename, free_cluster, FAT_EOF, fs_format)
	if err != nil {
		// fmt.Println("Error updating FAT entry for root directory:", err)
		return
	}

	// **Set the current and parent directory for the root directory**
	SetCurrentAndParentDirectory(filename, free_cluster, free_cluster, fs_format)

	// **Set the current cluster to the root directory**
	SetCurrentCluster(free_cluster)

	// fmt.Println("*** Root directory created successfully! ***")

}

func CreateDirectory(filename, dir_name string, fs_format FileSystemFormat) {

	// fmt.Println("*** Creating directory ***")

	// **Check if the directory name is valid**
	if dir_name == "." || dir_name == ".." {
		// fmt.Println("Error: Invalid directory name.")
		return
	}

	// **Parse the path to get the parent cluster and final directory name**
	parent_cluster, final_name, err := ParsePath(filename, dir_name, fs_format, true)
	if err != nil {
		// fmt.Println("Error parsing path:", err)
		fmt.Println("PATH NOT FOUND")
		return
	}

	// **Check if the directory name is too long**
	if len(final_name) > MAX_FILE_NAME {
		// fmt.Println("Error: Directory name is too long.")
		return
	}

	// **Check if the directory already exists**
	if CheckIfDirectoryExists(filename, parent_cluster, final_name, fs_format) {
		// fmt.Println("Error: Directory or file with the name", final_name, "already exists.")
		fmt.Println("EXIST")
		return
	}

	// **Find a free cluster for the new directory**
	free_cluster, err := FindFreeCluster(filename, fs_format.fat1_start)
	if err != nil {
		// fmt.Println("Error finding free cluster:", err)
		return
	}

	dir_name_bytes := [MAX_FILE_NAME]byte{}
	copy(dir_name_bytes[:], final_name)

	// **Create the new directory entry**
	new_dir := DirectoryEntry{
		Name:          dir_name_bytes,
		Size:          0,
		First_cluster: free_cluster,
		Is_directory:  1,
	}

	// **Write the directory entry to the file**
	// err = WriteDirectoryEntry(filename, free_cluster, new_dir, fs_format)
	// if err != nil {
	// 	// fmt.Println("Error writing directory entry:", err)
	// 	return
	// }

	// **Update the parent directory entry**
	err = UpdateParentDirectory(filename, parent_cluster, new_dir, fs_format)
	if err != nil {
		// fmt.Println("Error updating parent directory:", err)
		return
	}

	// **Update the FAT entry for the new directory**
	err = UpdateFatEntry(filename, free_cluster, FAT_EOF, fs_format)
	if err != nil {
		// fmt.Println("Error updating FAT entry:", err)
		return
	}

	// **Set the current and parent directory for the new directory**
	SetCurrentAndParentDirectory(filename, free_cluster, parent_cluster, fs_format)

	// fmt.Printf("Directory '%s' created at cluster %d.\n", final_name, free_cluster)
	// fmt.Println("*** Directory created successfully! ***")
}

func SetCurrentAndParentDirectory(filename string, current_cluster, parent_cluster int32, fs_format FileSystemFormat) {

	// fmt.Println("*** Setting current and parent directory ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// **Calculate the data cluster position for the directory entry**
	var offset int64
	if current_cluster == parent_cluster {
		offset = int64(fs_format.data_start)
	} else {
		data_cluster := current_cluster - 2*fs_format.fat_cluster_count - 1
		offset = int64(fs_format.data_start + data_cluster*CLUSTER_SIZE)
	}

	// **Current directory entry**
	current_entry := DirectoryEntry{
		Name:          [MAX_FILE_NAME]byte{'.'},
		Size:          0,
		First_cluster: current_cluster,
		Is_directory:  1,
	}

	if _, err := file.Seek(offset, 0); err != nil {
		// fmt.Println("Error seeking to '.' entry position:", err)
		return
	}
	if err := binary.Write(file, binary.LittleEndian, current_entry); err != nil {
		// fmt.Println("Error writing '.' entry:", err)
		return
	}

	// fmt.Println("Current cluster:", current_cluster)
	// fmt.Println("Parent cluster:", parent_cluster)
	// fmt.Println("Seeking to current directory position:", offset)
	// fmt.Println("Writing '.':", current_entry)
	// fmt.Println("Size of current entry:", binary.Size(current_entry))
	// fmt.Println()

	// **Parent directory entry**
	parent_entry := DirectoryEntry{
		Name:          [MAX_FILE_NAME]byte{'.', '.'},
		Size:          0,
		First_cluster: parent_cluster,
		Is_directory:  1,
	}

	if _, err := file.Seek(offset+int64(binary.Size(current_entry)), 0); err != nil {
		// fmt.Println("Error seeking to '..' entry position:", err)
		return
	}
	if err := binary.Write(file, binary.LittleEndian, parent_entry); err != nil {
		// fmt.Println("Error writing '..' entry:", err)
		return
	}

	// fmt.Println("Seeking to parent directory position:", offset+int64(binary.Size(current_entry)))
	// fmt.Println("Writing '..':", parent_entry)
	// fmt.Println("Size of parent entry:", binary.Size(parent_entry))
	// fmt.Println()

	// **Zero padding for the remaining space in the cluster**
	written_size := 2 * int64(binary.Size(current_entry))
	remaining_size := CLUSTER_SIZE - written_size
	zero_padding := make([]byte, remaining_size)

	// fmt.Println("Written size:", writtenSize)
	// fmt.Println("Remaining size:", remainingSize)
	// fmt.Println("Zero padding size:", len(zeroPadding))

	if _, err := file.Write(zero_padding); err != nil {
		// fmt.Println("Error padding remaining space with zeros:", err)
		return
	}

	// fmt.Println("Zero padding written successfully!")
	// fmt.Println("*** Current and parent directory set successfully! ***")
	// fmt.Println()
}

func CheckIfDirectoryExists(filename string, parent_cluster int32, dirName string, fs_format FileSystemFormat) bool {

	// fmt.Println("*** Checking if directory exists ***")

	// **Read the directory entries from the parent cluster**
	dir_entries, err := ReadDirectoryEntries(filename, parent_cluster, fs_format)
	if err != nil {
		// fmt.Println("Error reading directory entries:", err)
		return false
	}

	// **Check if the directory exists in the parent cluster**
	for _, entry := range dir_entries {

		if IsZeroEntry(entry) {
			continue
		}

		if string(bytes.Trim(entry.Name[:], "\x00")) == dirName {
			// fmt.Println("Directory exists!")
			return true
		}
	}

	return false
}

func ReadDirectoryEntries(filename string, cluster int32, fs_format FileSystemFormat) ([]DirectoryEntry, error) {

	// fmt.Println("*** Reading directory entries ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// **Calculate the data cluster position for the directory entry**
	offset := int64(fs_format.data_start + (cluster-2*fs_format.fat_cluster_count-1)*CLUSTER_SIZE)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking to cluster: %v", err)
	}

	// fmt.Println("Seeking to cluster:", cluster)
	// fmt.Println("Seeked to offset:", offset)

	// **Read the directory entries from the file**
	var items []DirectoryEntry
	for i := 0; i < CLUSTER_SIZE/binary.Size(DirectoryEntry{}); i++ {

		var entry DirectoryEntry
		err = binary.Read(file, binary.LittleEndian, &entry)
		if err != nil {
			return nil, fmt.Errorf("error reading directory entry: %v", err)
		}

		items = append(items, entry)

	}

	// fmt.Println("Directory entries read count:", len(items))
	// fmt.Println("*** Directory entries read successfully! ***")
	// fmt.Println()

	return items, nil
}

func IsZeroEntry(entry DirectoryEntry) bool {
	return entry.Name[0] == 0 && entry.Size == 0 && entry.First_cluster == 0
}

func UpdateFatEntry(filename string, cluster, value int32, fs_format FileSystemFormat) error {

	// fmt.Println("*** Updating FAT entry ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// **Seek to the FAT1 entry position**
	offset := int64(fs_format.fat1_start + cluster*FAT_ENTRY)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT entry: %v", err)
	}

	err = binary.Write(file, binary.LittleEndian, value)
	if err != nil {
		return fmt.Errorf("error updating FAT entry: %v", err)
	}
	// fmt.Println("Updating FAT1 entry at cluster", cluster)
	// fmt.Println("Seeked to offset:", offset)

	// **Seek to the FAT2 entry position**
	offset = int64(fs_format.fat2_start + cluster*FAT_ENTRY)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to FAT entry: %v", err)
	}

	err = binary.Write(file, binary.LittleEndian, value)
	if err != nil {
		return fmt.Errorf("error updating FAT entry: %v", err)
	}

	// fmt.Println("Updating FAT2 entry at cluster", cluster)
	// fmt.Println("Seeked to offset:", offset)

	// fmt.Println("*** FAT entry updated successfully! ***")
	// fmt.Println()

	return nil
}

func UpdateParentDirectory(filename string, parent_cluster int32, new_dir DirectoryEntry, fs_format FileSystemFormat) error {

	// fmt.Println("*** Updating parent directory ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// **Read the directory entries from the parent cluster**
	dir_entries, err := ReadDirectoryEntries(filename, parent_cluster, fs_format)
	if err != nil {
		return fmt.Errorf("error reading directory entries: %v", err)
	}

	// **Find the free entry in the parent directory**
	entry_written := false
	for i, entry := range dir_entries {
		if IsZeroEntry(entry) {
			dir_entries[i] = new_dir
			entry_written = true
			// fmt.Println("Free entry found in parent directory:", i)
			break
		}
	}

	// **If no free entry was found, find a new cluster for the parent directory**
	if !entry_written {

		// fmt.Println("No free entry found in parent directory. Finding a new cluster...")

		// **Find a new free cluster**
		new_cluster, err := FindFreeCluster(filename, fs_format.fat1_start)
		if err != nil {
			return fmt.Errorf("error finding free cluster: %v", err)
		}

		// **Update the parent directory's FAT entry to link to the new cluster**
		err = UpdateFatEntry(filename, parent_cluster, new_cluster, fs_format)
		if err != nil {
			return fmt.Errorf("error updating FAT entry for parent directory: %v", err)
		}

		// **Write the new directory entries to the new cluster**
		err = WriteDirectoryEntry(filename, new_cluster, new_dir, fs_format)
		if err != nil {
			return fmt.Errorf("error writing directory entries to new cluster: %v", err)
		}

		// **Update the parent directory entry to link to the new cluster**
		err = WriteDirectoryEntry(filename, parent_cluster, new_dir, fs_format)
		if err != nil {
			return fmt.Errorf("error writing new directory entry to parent directory: %v", err)
		}
	}

	// **Write the updated directory entries back to the parent cluster**
	// fmt.Println("Writing updated directory entries to parent directory...")

	// **Seek to the parent directory cluster position**
	offset := int64(fs_format.data_start + (parent_cluster-2*fs_format.fat_cluster_count-1)*CLUSTER_SIZE)
	// fmt.Println("Seeking to parent directory cluster position:", offset)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to parent directory cluster: %v", err)
	}

	// fmt.Println("Offset:", offset)

	// **Write the directory entries to the file**
	for _, entry := range dir_entries {

		if IsZeroEntry(entry) {
			continue
		}

		// fmt.Println("Writing directory entry:", entry)
		err = binary.Write(file, binary.LittleEndian, entry)
		if err != nil {
			return fmt.Errorf("error writing directory entry: %v", err)
		}
	}

	// fmt.Println("*** Parent directory updated successfully! ***")
	// fmt.Println()
	return nil
}

func GetCurrentCluster() int32 {
	// fmt.Println("Getting current cluster:", current_cluster)
	return current_cluster
}

func SetCurrentCluster(cluster int32) {
	current_cluster = cluster
	// fmt.Println("Setting current cluster:", current_cluster)
}

func GetParentCluster(filename string, current_cluster int32, fs_format FileSystemFormat) int32 {

	// fmt.Println("*** Getting parent cluster ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return -1
	}

	// **Calculate the data cluster position for the directory entry**
	data_cluster := current_cluster - 2*fs_format.fat_cluster_count - 1
	offset := int64(fs_format.data_start + data_cluster*CLUSTER_SIZE)

	// **Seek to the parent directory entry position**
	_, err = file.Seek(offset+int64(binary.Size(DirectoryEntry{})), 0)
	if err != nil {
		// fmt.Println("Error seeking to parent directory entry:", err)
		return -1
	}

	// **Read the parent directory entry from the file**
	var parent_entry DirectoryEntry
	err = binary.Read(file, binary.LittleEndian, &parent_entry)
	if err != nil {
		// fmt.Println("Error reading parent directory entry:", err)
		return -1
	}

	// fmt.Println("Parent directory entry:", parent_entry)
	// fmt.Println("Parent cluster:", parent_entry.First_cluster)
	// fmt.Println("*** Parent cluster retrieved successfully! ***")
	// fmt.Println()

	return parent_entry.First_cluster
}

func GetCurrentPath() string {
	fmt.Println("Getting current path:", current_path)
	return current_path
}

func SetCurrentPath(path string) {

	if path == "/" {
		current_path = "/"
	} else if path == "." {
		return
	} else if path == ".." {

		if current_path == "/" {
			return
		}

		// Remove the last directory from the current path
		current_path = current_path[:len(current_path)-1]
		for i := len(current_path) - 1; i >= 0; i-- {
			if current_path[i] == '/' {
				current_path = current_path[:i+1]
				break
			}
		}

	} else {
		current_path += path + "/"
	}
	// fmt.Println("Setting current path:", current_path)
}

func ParsePath(filename, dest string, fs_format FileSystemFormat, last_entry bool) (int32, string, error) {
	// fmt.Println("*** Parsing path ***")

	// **Trim the trailing slash from the destination path**
	dest = strings.TrimRight(dest, "/")
	path_components := strings.Split(dest, "/")

	// **Check if the path is absolute or relative**
	var current_cluster int32
	if path_components[0] == "" {
		// Absolute path (starts with "/"): Start from the root directory
		current_cluster = fs_format.data_start / CLUSTER_SIZE
	} else {
		// Relative path: Start from the current directory
		current_cluster = GetCurrentCluster()
	}

	// **Split the path into components**
	var final_name string
	for i, component := range path_components {

		if component == "" || component == "." {
			continue // Ignore empty or current directory symbol
		}

		if component == ".." {
			// Handle parent directory navigation
			parent_cluster := GetParentCluster(filename, current_cluster, fs_format)
			current_cluster = parent_cluster
			continue
		}

		if i == len(path_components)-1 {
			final_name = component

			// If `last_entry` is false, traverse into the last component
			if !last_entry {
				next_cluster, err := FindDirectoryCluster(filename, component, current_cluster, fs_format)
				if err != nil {
					fmt.Println("Error finding cluster for last component:", err)
					return -1, "", fmt.Errorf("error finding cluster for '%s': %v", component, err)
				}
				if next_cluster == -1 {
					fmt.Println("Directory or file not found:", component)
					return -1, "", fmt.Errorf("directory or file '%s' not found", component)
				}
				current_cluster = next_cluster
			}
			break
		}

		// Traverse to the next directory
		next_cluster, err := FindDirectoryCluster(filename, component, current_cluster, fs_format)
		if err != nil {
			return -1, "", fmt.Errorf("error finding cluster for directory '%s': %v", component, err)
		}
		if next_cluster == -1 {
			return -1, "", fmt.Errorf("directory '%s' not found in path", component)
		}
		current_cluster = next_cluster
	}

	// fmt.Println("*** Path parsed successfully! ***")
	return current_cluster, final_name, nil
}

func RemoveDirectoryEntry(filename string, cluster int32, dir_name string, fs_format FileSystemFormat) error {

	if dir_name == "." || dir_name == ".." || dir_name == "/" {
		// fmt.Println("Error: Invalid directory name.")
		return nil
	}

	// fmt.Println("*** Removing directory entry ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// **Read the directory entries from the cluster**
	dir_entries, err := ReadDirectoryEntries(filename, cluster, fs_format)
	if err != nil {
		return fmt.Errorf("error reading directory entries: %v", err)
	}

	// **Find the directory entry to remove**
	entry_index := -1
	for i, entry := range dir_entries {
		if string(bytes.Trim(entry.Name[:], "\x00")) == dir_name {
			entry_index = i
			// fmt.Println("Directory entry found at index:", i)
			break
		}
	}

	if entry_index == -1 {
		// fmt.Println("Error: Directory entry not found.")
		fmt.Println("FILE NOT FOUND")
		return nil
	}

	// **Check if the directory has contents and prevent removal if not empty**
	entry_to_remove := dir_entries[entry_index]

	if entry_to_remove.Is_directory == 1 {

		// fmt.Println("Directory entry to remove:", entry_to_remove)
		sub_entries, err := ReadDirectoryEntries(filename, entry_to_remove.First_cluster, fs_format)
		if err != nil {
			return fmt.Errorf("error reading subdirectory entries: %v", err)
		}

		// Check if the subdirectory is empty
		for _, sub_entry := range sub_entries {

			name := string(bytes.Trim(sub_entry.Name[:], "\x00"))

			if name == "." || name == ".." {
				continue
			}

			if !IsZeroEntry(sub_entry) {
				fmt.Println("NOT EMPTY")
				return fmt.Errorf("directory '%s' is not empty", dir_name)
			}
		}
	}

	// **Clear the FAT entries for the directory's clusters**
	cluster_to_clear := entry_to_remove.First_cluster
	for cluster_to_clear != FAT_EOF {

		// fmt.Println("Clearing cluster:", cluster_to_clear)

		next_cluster, err := ReadFatEntry(filename, cluster_to_clear, fs_format)
		if err != nil {
			return fmt.Errorf("error reading FAT entry: %v", err)
		}

		// fmt.Println("Next cluster:", next_cluster)

		// Mark the current cluster as free
		err = UpdateFatEntry(filename, cluster_to_clear, FAT_FREE, fs_format)
		if err != nil {
			return fmt.Errorf("error clearing FAT entry: %v", err)
		}

		cluster_to_clear = next_cluster
	}

	// **Remove the directory entry by clearing it**
	dir_entries[entry_index] = DirectoryEntry{}
	// fmt.Println("Directory entry removed:", dir_name)

	// **Write the updated directory entries back to the cluster**
	offset := int64(fs_format.data_start + (cluster-2*fs_format.fat_cluster_count-1)*CLUSTER_SIZE)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to cluster: %v", err)
	}

	for _, entry := range dir_entries {
		err = binary.Write(file, binary.LittleEndian, entry)
		if err != nil {
			return fmt.Errorf("error writing directory entry: %v", err)
		}
	}

	// fmt.Println("*** Directory entry and its contents removed successfully! ***")
	return nil
}

func ReadFatEntry(filename string, cluster int32, fs_format FileSystemFormat) (int32, error) {

	// Calculate the offset in the FAT table for the given cluster
	offset := int64(fs_format.fat1_start + cluster*FAT_ENTRY)

	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Seek to the position of the FAT entry
	_, err = file.Seek(offset, 0)
	if err != nil {
		return 0, fmt.Errorf("error seeking to FAT entry: %v", err)
	}

	// Read the FAT entry
	var nextCluster int32
	err = binary.Read(file, binary.LittleEndian, &nextCluster)
	if err != nil {
		return 0, fmt.Errorf("error reading FAT entry: %v", err)
	}

	return nextCluster, nil
}

func FindDirectoryCluster(filename, dir_name string, parent_cluster int32, fs_format FileSystemFormat) (int32, error) {

	// Read the directory entries from the parent cluster
	dir_entries, err := ReadDirectoryEntries(filename, parent_cluster, fs_format)
	if err != nil {
		return -1, fmt.Errorf("error reading directory entries: %v", err)
	}

	// Find the directory entry in the parent cluster
	for _, entry := range dir_entries {
		if string(bytes.Trim(entry.Name[:], "\x00")) == dir_name {
			return entry.First_cluster, nil
		}
	}

	return -1, nil
}

func FindEntry(filename, src string, current_cluster int32, fs_format FileSystemFormat) (DirectoryEntry, error) {

	// fmt.Println("*** Checking file ***")

	// **Read the directory entries from the cluster**
	dir_entries, err := ReadDirectoryEntries(filename, current_cluster, fs_format)
	if err != nil {
		return DirectoryEntry{}, fmt.Errorf("error reading directory entries: %v", err)
	}

	// **Find the file entry in the cluster**
	for _, entry := range dir_entries {
		if string(bytes.Trim(entry.Name[:], "\x00")) == src {
			return entry, nil
		}
	}

	return DirectoryEntry{}, fmt.Errorf("error reading directory entries: %v", err)
}

func ReadCluster(filename string, cluster int32, fs_format FileSystemFormat) error {

	// fmt.Println("*** Reading cluster ***")

	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// **Calculate the data cluster position for the directory entry**
	offset := int64(fs_format.data_start + (cluster-2*fs_format.fat_cluster_count-1)*CLUSTER_SIZE)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("error seeking to cluster: %v", err)
	}

	// **Read the cluster from the file**
	cluster_data := make([]byte, CLUSTER_SIZE)
	_, err = file.Read(cluster_data)
	if err != nil {
		return fmt.Errorf("error reading cluster: %v", err)
	}

	// **Print the cluster data**
	fmt.Println("Cluster data:")
	fmt.Println(string(cluster_data))

	// fmt.Println("Cluster read successfully!")
	// fmt.Println("*** Cluster read successfully! ***")
	// fmt.Println()

	return nil
}

func ReadFileContents(filename string, start_cluster int32, file_size int32, fs_format FileSystemFormat) ([]byte, error) {

	var file_contents []byte
	current_cluster := start_cluster
	remaining_size := file_size

	// Open the VFS file
	vfsFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening VFS file: %v", err)
	}
	defer vfsFile.Close()

	for remaining_size > 0 {
		// Calculate the offset for the current cluster
		offset := int64(fs_format.data_start + (current_cluster-2*fs_format.fat_cluster_count-1)*CLUSTER_SIZE)
		readSize := CLUSTER_SIZE
		if remaining_size < CLUSTER_SIZE {
			readSize = int(remaining_size)
		}

		// Read the cluster's data
		buffer := make([]byte, readSize)
		bytesRead, err := vfsFile.ReadAt(buffer, offset)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading cluster %d at offset %d: %v", current_cluster, offset, err)
		}

		// Append the data to the file_contents
		file_contents = append(file_contents, buffer...)

		// Reduce the remaining size
		remaining_size -= int32(bytesRead)

		// Stop reading if EOF reached or remaining size is zero
		if remaining_size <= 0 {
			break
		}

		// Get the next cluster from FAT
		current_cluster, err = ReadFatEntry(filename, current_cluster, fs_format)
		if err != nil {
			return nil, fmt.Errorf("error reading FAT entry for cluster %d: %v", current_cluster, err)
		}

		// Check if we've reached the end of the file
		if current_cluster == FAT_EOF {
			break
		}

	}

	return file_contents, nil
}

func WriteFileContents(filename string, startCluster int32, file_contents []byte, fs_format FileSystemFormat) error {

	// fmt.Println("*** Writing file contents ***")

	current_cluster := startCluster
	remaining_size := int32(len(file_contents))

	// fmt.Println("Start cluster:", startCluster)
	// fmt.Println("Remaining size:", remaining_size)

	// Open the VFS file
	vfsFile, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening VFS file: %v", err)
	}
	defer vfsFile.Close()

	for remaining_size > 0 {
		// Calculate the offset for the current cluster
		offset := int64(fs_format.data_start + (current_cluster-2*fs_format.fat_cluster_count-1)*CLUSTER_SIZE)
		writeSize := CLUSTER_SIZE
		if remaining_size < CLUSTER_SIZE {
			writeSize = int(remaining_size)
		}

		// fmt.Println("Writing to cluster:", current_cluster)
		// fmt.Println("Offset:", offset)

		// Write the cluster's data
		_, err := vfsFile.WriteAt(file_contents[:writeSize], offset)
		if err != nil {
			return fmt.Errorf("error writing cluster %d: %v", current_cluster, err)
		}

		// fmt.Println("Data written to cluster:", current_cluster)

		// Update the current cluster and remaining size
		err = UpdateFatEntry(filename, current_cluster, FAT_EOF, fs_format)
		if err != nil {
			return fmt.Errorf("error updating FAT entry for cluster %d: %v", current_cluster, err)
		}

		// fmt.Println("Updated FAT entry for cluster:", current_cluster, "to:", FAT_EOF)

		file_contents = file_contents[writeSize:]
		remaining_size -= int32(writeSize)
		// fmt.Println("Remaining size to write:", remaining_size)

		if remaining_size <= 0 {
			break
		}

		// Get the next cluster from FAT
		nextCluster, err := FindFreeCluster(filename, fs_format.fat1_start)
		if err != nil {
			return fmt.Errorf("error finding free cluster: %v", err)
		}

		// Update the current cluster and remaining size
		// fmt.Println("Next cluster:", nextCluster)

		// Update the FAT entry for the current cluster
		err = UpdateFatEntry(filename, current_cluster, nextCluster, fs_format)
		if err != nil {
			return fmt.Errorf("error updating FAT entry for cluster %d: %v", current_cluster, err)
		}

		// fmt.Println("Updated FAT entry for cluster:", current_cluster, "to:", nextCluster)

		current_cluster = nextCluster
	}

	// fmt.Println("*** File contents written successfully! ***")
	return nil
}
