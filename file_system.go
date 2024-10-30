package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

// Initializes the file system
func (fs *FileSystem) Init() {

	fs.fat_table = FAT{}

	for i := range fs.fat_table {
		fs.fat_table[i] = FAT_FREE
	}

	fs.directory = make(map[string]DirectoryEntry)
	fs.cluster_data = make([][]byte, MAX_CLUSTER_COUNT)

	for i := range fs.cluster_data {
		fs.cluster_data[i] = make([]byte, CLUSTER_SIZE)
	}

}

// Find a free cluster in FAT
func (fs *FileSystem) FindFreeCluster() (int, error) {

	for i, val := range fs.fat_table {
		if val == FAT_FREE {
			return i, nil
		}
	}

	return -1, errors.New("no free clusters available")
}

// Load file from external system into pseudo-FAT (incp s1 s2)
func (fs *FileSystem) InCp(filename string, data []byte) error {

	// Ensure the filename is not too long
	if len(filename) > MAX_FILE_NAME {
		return fmt.Errorf("filename too long")
	}

	// Ensure the file does not already exist
	if _, exists := fs.directory[filename]; exists {
		return fmt.Errorf("file already exists")
	}

	// Split data into clusters
	numClusters := (len(data) + CLUSTER_SIZE - 1) / CLUSTER_SIZE
	firstCluster, err := fs.FindFreeCluster()
	if err != nil {
		return err
	}

	currentCluster := firstCluster
	for i := 0; i < numClusters; i++ {
		if i > 0 {
			newCluster, err := fs.FindFreeCluster()
			if err != nil {
				return err
			}
			fs.fat_table[currentCluster] = newCluster
			currentCluster = newCluster
		}

		// Copy the data into the cluster
		start := i * CLUSTER_SIZE
		end := start + CLUSTER_SIZE
		if end > len(data) {
			end = len(data)
		}
		copy(fs.cluster_data[currentCluster], data[start:end])
	}
	fs.fat_table[currentCluster] = FAT_EOF

	// Add to directory
	fs.directory[filename] = DirectoryEntry{
		name:          filename,
		size:          len(data),
		first_cluster: firstCluster,
	}

	fmt.Println("OK")
	return nil
}

// Read a file from the pseudo-FAT (cat s1)
func (fs *FileSystem) Cat(filename string) error {

	entry, exists := fs.directory[filename]
	if !exists {
		return fmt.Errorf("FILE NOT FOUND")
	}

	currentCluster := entry.first_cluster
	for currentCluster != FAT_EOF {
		fmt.Printf("%s", fs.cluster_data[currentCluster])
		currentCluster = fs.fat_table[currentCluster]
	}
	fmt.Println("OK")
	return nil
}

// Remove a file (rm s1)
func (fs *FileSystem) Rm(filename string) error {
	entry, exists := fs.directory[filename]
	if !exists {
		return fmt.Errorf("FILE NOT FOUND")
	}

	currentCluster := entry.first_cluster
	for currentCluster != FAT_EOF {
		nextCluster := fs.fat_table[currentCluster]
		fs.fat_table[currentCluster] = FAT_FREE // Mark the cluster as free
		currentCluster = nextCluster
	}

	delete(fs.directory, filename)
	fmt.Println("OK")
	return nil
}

func SaveFileSystem(filename string, fs *FileSystem, fs_format FileSystemFormat) error {

	fmt.Printf("\nSaving file system to '%s'...\n", filename)

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
	for _, val := range fs.fat_table {
		err = binary.Write(file, binary.LittleEndian, int32(val))
		if err != nil {
			return fmt.Errorf("error writing FAT1 table: %w", err)
		}
	}
	fmt.Println("FAT1 table written successfully")

	// **Write FAT2 table at fat2_start position**
	fmt.Println("Seeking to FAT2 start position:", fs_format.fat2_start)
	_, err = file.Seek(int64(fs_format.fat2_start), 0)

	if err != nil {
		return fmt.Errorf("error seeking to FAT2 start: %w", err)
	}

	for _, val := range fs.fat_table {
		err = binary.Write(file, binary.LittleEndian, int32(val))
		if err != nil {
			return fmt.Errorf("error writing FAT2 table: %w", err)
		}
	}

	fmt.Println("FAT2 table written successfully")

	// **Zero out the data starting at data_start**
	fmt.Println("Seeking to data start position:", fs_format.data_start)
	_, err = file.Seek(int64(fs_format.data_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to data start: %w", err)
	}

	// **Zero out the data starting at data_start**
	fmt.Println("Zeroing out data section...")
	remainingSize := fs_format.file_size - fs_format.data_start
	zeroBuffer := make([]byte, remainingSize) // Create a buffer filled with zeros
	fmt.Println("Zero buffer size:", len(zeroBuffer))
	fmt.Println("Seeking to data start position:", fs_format.data_start)
	_, err = file.Seek(int64(fs_format.data_start), 0)
	if err != nil {
		return fmt.Errorf("error seeking to data start: %w", err)
	}
	_, err = file.Write(zeroBuffer)
	if err != nil {
		return fmt.Errorf("error writing zeros to data section: %w", err)
	}

	fmt.Println("Data section zeroed out successfully")
	fmt.Println("File system saved successfully!")
	return nil
}

func LoadFileSystem(filename string) *FileSystem {
	fmt.Printf("\nLoading file system from '%s'...\n", filename)

	// Load the format file first to get the necessary information
	fs_format := LoadFormatFile(filename)
	if fs_format.file_size == 0 {
		fmt.Println("Format file is empty or could not be read.")
		return nil
	}

	// Now we can open the file and read the file system data
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	fs := &FileSystem{}
	fs.Init()

	// Read the FAT table
	fmt.Printf("Reading FAT table of size %d bytes...\n", fs_format.fat_size)
	_, err = file.Seek(int64(fs_format.fat1_start), 0) // Seek to FAT1 start
	if err != nil {
		fmt.Println("Error seeking to FAT1 start:", err)
		return nil
	}
	for i := range fs.fat_table {
		var val int32
		err = binary.Read(file, binary.LittleEndian, &val)
		if err != nil {
			fmt.Println("Error reading FAT table:", err)
			return nil
		}
		fs.fat_table[i] = int(val)
	}
	fmt.Println("FAT table loaded successfully.")

	// Read the directory (Assuming a structure for your directory)
	fmt.Println("Reading directory entries...")
	for {
		var name string
		for {
			var b byte
			_, err = file.Read([]byte{b})
			if err != nil {
				fmt.Println("Error reading file:", err)
				return nil
			}
			if b == 0 { // Null terminator indicates end of name
				break
			}
			name += string(b)
		}
		if name == "" {
			break // Exit if no name found
		}

		var entry DirectoryEntry
		err = binary.Read(file, binary.LittleEndian, &entry.size)
		if err != nil {
			fmt.Println("Error reading entry size:", err)
			return nil
		}
		err = binary.Read(file, binary.LittleEndian, &entry.first_cluster)
		if err != nil {
			fmt.Println("Error reading entry first cluster:", err)
			return nil
		}
		fs.directory[name] = entry
	}
	fmt.Println("Directory loaded successfully.")

	// Read the cluster data based on total cluster count
	fmt.Printf("Reading cluster data...\n")
	for i := 0; i < int(fs_format.cluster_count); i++ {
		_, err = file.Read(fs.cluster_data[i])
		if err != nil {
			fmt.Println("Error reading cluster data:", err)
			return nil
		}
	}
	fmt.Println("Cluster data loaded successfully.")

	fmt.Println("File system loaded successfully!")
	return fs
}

func SaveFormatFile(filename string, fs_format FileSystemFormat) {
	fmt.Printf("\nSaving format of the file system to '%s'...\n", filename)

	// Open the file for writing
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the file size
	err = binary.Write(file, binary.LittleEndian, int32(fs_format.file_size))
	if err != nil {
		fmt.Println("Error writing file size:", err)
		return
	}
	fmt.Printf("File size written: %d bytes\n", fs_format.file_size)

	// Write the FAT size
	err = binary.Write(file, binary.LittleEndian, int32(fs_format.fat_size))
	if err != nil {
		fmt.Println("Error writing FAT size:", err)
		return
	}
	fmt.Printf("FAT size written: %d bytes\n", fs_format.fat_size)

	// Write the number of FAT clusters
	err = binary.Write(file, binary.LittleEndian, int32(fs_format.fat_cluster_count))
	if err != nil {
		fmt.Println("Error writing FAT cluster count:", err)
		return
	}
	fmt.Printf("FAT cluster count written: %d\n", fs_format.fat_cluster_count)

	// Write the total number of data clusters
	err = binary.Write(file, binary.LittleEndian, int32(fs_format.cluster_count))
	if err != nil {
		fmt.Println("Error writing cluster count:", err)
		return
	}
	fmt.Printf("Cluster count written: %d\n", fs_format.cluster_count)

	// Write the starting positions
	err = binary.Write(file, binary.LittleEndian, int32(fs_format.fat1_start))
	if err != nil {
		fmt.Println("Error writing FAT1 start position:", err)
		return
	}
	fmt.Printf("FAT1 starts at: %d\n", fs_format.fat1_start)

	err = binary.Write(file, binary.LittleEndian, int32(fs_format.fat2_start))
	if err != nil {
		fmt.Println("Error writing FAT2 start position:", err)
		return
	}
	fmt.Printf("FAT2 starts at: %d\n", fs_format.fat2_start)

	err = binary.Write(file, binary.LittleEndian, int32(fs_format.data_start))
	if err != nil {
		fmt.Println("Error writing data start position:", err)
		return
	}
	fmt.Printf("Data starts at: %d\n", fs_format.data_start)

	fmt.Println("File system format saved successfully!")
}

func LoadFormatFile(filename string) FileSystemFormat {
	fmt.Printf("\nLoading format of the file system from '%s'...\n", filename)

	// Initialize the FileSystemFormat struct
	fs_format := FileSystemFormat{}

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return FileSystemFormat{}
	}
	defer file.Close()

	// Read the file size
	err = binary.Read(file, binary.LittleEndian, &fs_format.file_size)
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return FileSystemFormat{}
	}
	fmt.Printf("File size read: %d bytes\n", fs_format.file_size)

	// Read the FAT size
	err = binary.Read(file, binary.LittleEndian, &fs_format.fat_size)
	if err != nil {
		fmt.Println("Error reading FAT size:", err)
		return FileSystemFormat{}
	}
	fmt.Printf("FAT size read: %d bytes\n", fs_format.fat_size)

	// Read the number of FAT clusters
	err = binary.Read(file, binary.LittleEndian, &fs_format.fat_cluster_count)
	if err != nil {
		fmt.Println("Error reading FAT cluster count:", err)
		return FileSystemFormat{}
	}
	fmt.Printf("FAT cluster count read: %d\n", fs_format.fat_cluster_count)

	// Read the total number of data clusters
	err = binary.Read(file, binary.LittleEndian, &fs_format.cluster_count)
	if err != nil {
		fmt.Println("Error reading cluster count:", err)
		return FileSystemFormat{}
	}
	fmt.Printf("Cluster count read: %d\n", fs_format.cluster_count)

	// Read the starting positions
	err = binary.Read(file, binary.LittleEndian, &fs_format.fat1_start)
	if err != nil {
		fmt.Println("Error reading FAT1 start position:", err)
		return FileSystemFormat{}
	}
	fmt.Printf("FAT1 starts at: %d\n", fs_format.fat1_start)

	err = binary.Read(file, binary.LittleEndian, &fs_format.fat2_start)
	if err != nil {
		fmt.Println("Error reading FAT2 start position:", err)
		return FileSystemFormat{}
	}
	fmt.Printf("FAT2 starts at: %d\n", fs_format.fat2_start)

	err = binary.Read(file, binary.LittleEndian, &fs_format.data_start)
	if err != nil {
		fmt.Println("Error reading data start position:", err)
		return FileSystemFormat{}
	}
	fmt.Printf("Data starts at: %d\n", fs_format.data_start)

	fmt.Println("File system format loaded successfully!")
	return fs_format
}

// Format the file with the desired size
func FormatFile(filename string, file_size int) {

	// Initialize the file system in memory (FAT, directory, etc.)
	fs := &FileSystem{}

	// Adjust the number of clusters based on the file size
	cluster_count := int(file_size / CLUSTER_SIZE)
	if cluster_count > MAX_CLUSTER_COUNT {
		cluster_count = MAX_CLUSTER_COUNT
	}

	fmt.Printf("Formatting file system with %d clusters...\n", cluster_count)

	// Calculate the size of the FAT in bytes and the number of clusters it occupies
	fat_size := cluster_count * FAT_ENTRY
	fat_cluster_count := (fat_size + CLUSTER_SIZE - 1) / CLUSTER_SIZE

	fmt.Printf("FAT size: %d bytes\nFAT clusters: %d\n", fat_size, fat_cluster_count)

	// TODO: Change the init function to accept the number of clusters
	// Initialize the file system's FAT and cluster data for the given size
	fs.fat_table = make(FAT, cluster_count)
	fs.cluster_data = make([][]byte, cluster_count)

	for i := 0; i < cluster_count; i++ {
		fs.fat_table[i] = FAT_FREE
		fs.cluster_data[i] = make([]byte, CLUSTER_SIZE)
	}

	// Initialize the directory
	fs.directory = make(map[string]DirectoryEntry)
	// TODO: move to the init, to this point

	fmt.Println()
	fmt.Println("FAT table size:", len(fs.fat_table))
	fmt.Println("Cluster data size:", len(fs.cluster_data))
	fmt.Println("Directory size:", len(fs.directory))
	fmt.Println()

	// Set the starting addresses for FAT1, FAT2, and data section (optional)
	fat1_start := CLUSTER_SIZE                                // FAT1 starts after boot sector
	fat2_start := fat1_start + fat_cluster_count*CLUSTER_SIZE // FAT2 starts after FAT1
	data_start := fat2_start + fat_cluster_count*CLUSTER_SIZE // Data starts after FAT2

	fmt.Printf("FAT1 starts at: %d\nFAT2 starts at: %d\nData starts at: %d\n", fat1_start, fat2_start, data_start)

	fs_format := FileSystemFormat{
		file_size:         int32(file_size),
		fat_size:          int32(fat_size),
		fat_cluster_count: int32(fat_cluster_count),
		cluster_count:     int32(cluster_count),
		fat1_start:        int32(fat1_start),
		fat2_start:        int32(fat2_start),
		data_start:        int32(data_start),
	}

	SaveFormatFile(filename, fs_format)

	// Save the file system to the file
	err := SaveFileSystem(filename, fs, fs_format)
	if err != nil {
		fmt.Println("Error saving file system:", err)
		return
	}

	fmt.Println("File system formatted and saved successfully!")
	fmt.Println()
}
