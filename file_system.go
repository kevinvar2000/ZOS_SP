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

func SaveFileSystem(filename string, fs *FileSystem) {

	fmt.Printf("\nSaving file system...\n\n")

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the FAT table
	for _, val := range fs.fat_table {
		err = binary.Write(file, binary.LittleEndian, int32(val))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Println("FAT table written successfully, size:", len(fs.fat_table))

	// Write the directory
	for name, entry := range fs.directory {

		// Write the file name as a fixed 12-byte field (8.3 filename format)
		nameBytes := make([]byte, MAX_FILE_NAME)
		copy(nameBytes, name) // Copy the file name to the fixed-length byte slice

		_, err = file.Write(nameBytes)
		if err != nil {
			return
		}

		fmt.Println("Writing file name with size:", len(nameBytes))

		// Write the file size as a 4-byte integer
		err = binary.Write(file, binary.LittleEndian, int32(entry.size))
		if err != nil {
			return
		}

		fmt.Println("Writing entry size with size:", entry.size)

		// Write the first cluster as a 4-byte integer
		err = binary.Write(file, binary.LittleEndian, int32(entry.first_cluster))
		if err != nil {
			return
		}

		fmt.Println("Writing first cluster with size:", entry.first_cluster)

	}

	// Write the cluster data
	// for _, data := range fs.cluster_data {
	// 	_, err = file.Write(data)
	// 	if err != nil {
	// 		fmt.Println("Error writing to file:", err)
	// 		return
	// 	}
	// }

	// Write the cluster data for the used clusters only
	// for i := 0; i < cluster_count; i++ { // Use the correct number of clusters
	// 	_, err = file.Write(fs.cluster_data[i])
	// 	if err != nil {
	// 		fmt.Println("Error writing to file:", err)
	// 		return
	// 	}
	// }

	// Write the cluster data for the used clusters only
	for _, entry := range fs.directory {
		numClusters := (entry.size + CLUSTER_SIZE - 1) / CLUSTER_SIZE // Calculate the total number of clusters needed for this file
		for i := 0; i < numClusters; i++ {
			// Determine the size of data to write for the last cluster
			var bytesToWrite int
			if i == numClusters-1 { // If this is the last cluster
				bytesToWrite = entry.size % CLUSTER_SIZE
				if bytesToWrite == 0 && i > 0 { // If perfectly divisible, write full last cluster
					bytesToWrite = CLUSTER_SIZE
				}
			} else {
				bytesToWrite = CLUSTER_SIZE // Full cluster for non-last clusters
			}

			// Write the data for the current cluster
			clusterIndex := entry.first_cluster + i
			if clusterIndex < len(fs.cluster_data) {
				_, err = file.Write(fs.cluster_data[clusterIndex][:bytesToWrite])
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
		}
	}

	fmt.Println("File system saved successfully!")

}

func LoadFileSystem(filename string) *FileSystem {

	fmt.Printf("\nLoading file system...\n\n")

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Check file size
	_, err = file.Stat()
	if err != nil {
		fmt.Println("Error file stat:", err)
		return nil
	}

	fs := &FileSystem{}
	fs.Init()

	// Read the FAT table
	for i := range fs.fat_table {
		var val byte
		_, err = file.Read([]byte{val})
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		fs.fat_table[i] = int(val)
	}

	// Read the directory
	for {
		var name string
		for {
			var b byte
			_, err = file.Read([]byte{b})
			if err != nil {
				fmt.Println("Error reading file:", err)
				return nil
			}
			if b == 0 {
				break
			}
			name += string(b)
		}
		if name == "" {
			break
		}

		var entry DirectoryEntry
		entry.name = name
		_, err = file.Read([]byte{byte(entry.size)})
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		_, err = file.Read([]byte{byte(entry.first_cluster)})
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		fs.directory[name] = entry
	}

	// Read the cluster data
	for i := range fs.cluster_data {
		_, err = file.Read(fs.cluster_data[i])
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
	}

	fmt.Println("File system loaded successfully!")
	return fs
}

func SaveFormatFile(filename string, file_size int, fat_size int, fat_cluster_count int, cluster_count int, fat1_start int, fat2_start int, data_start int) {
	fmt.Printf("\nSaving format of the file system to '%s'...\n", filename)

	// Open the file for writing
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the file size
	err = binary.Write(file, binary.LittleEndian, int32(file_size))
	if err != nil {
		fmt.Println("Error writing file size:", err)
		return
	}
	fmt.Printf("File size written: %d bytes\n", file_size)

	// Write the FAT size
	err = binary.Write(file, binary.LittleEndian, int32(fat_size))
	if err != nil {
		fmt.Println("Error writing FAT size:", err)
		return
	}
	fmt.Printf("FAT size written: %d bytes\n", fat_size)

	// Write the number of FAT clusters
	err = binary.Write(file, binary.LittleEndian, int32(fat_cluster_count))
	if err != nil {
		fmt.Println("Error writing FAT cluster count:", err)
		return
	}
	fmt.Printf("FAT cluster count written: %d\n", fat_cluster_count)

	// Write the total number of data clusters
	err = binary.Write(file, binary.LittleEndian, int32(cluster_count))
	if err != nil {
		fmt.Println("Error writing cluster count:", err)
		return
	}
	fmt.Printf("Cluster count written: %d\n", cluster_count)

	// Write the starting positions
	err = binary.Write(file, binary.LittleEndian, int32(fat1_start))
	if err != nil {
		fmt.Println("Error writing FAT1 start position:", err)
		return
	}
	fmt.Printf("FAT1 starts at: %d\n", fat1_start)

	err = binary.Write(file, binary.LittleEndian, int32(fat2_start))
	if err != nil {
		fmt.Println("Error writing FAT2 start position:", err)
		return
	}
	fmt.Printf("FAT2 starts at: %d\n", fat2_start)

	err = binary.Write(file, binary.LittleEndian, int32(data_start))
	if err != nil {
		fmt.Println("Error writing data start position:", err)
		return
	}
	fmt.Printf("Data starts at: %d\n", data_start)

	fmt.Println("File system format saved successfully!")
}

func LoadFormatFile(filename string) (int, int, int, int, int, int, int) {
	fmt.Printf("\nLoading format of the file system from '%s'...\n", filename)

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	defer file.Close()

	// Read the file size
	var file_size int32
	err = binary.Read(file, binary.LittleEndian, &file_size)
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("File size read: %d bytes\n", file_size)

	// Read the FAT size
	var fat_size int32
	err = binary.Read(file, binary.LittleEndian, &fat_size)
	if err != nil {
		fmt.Println("Error reading FAT size:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("FAT size read: %d bytes\n", fat_size)

	// Read the number of FAT clusters
	var fat_cluster_count int32
	err = binary.Read(file, binary.LittleEndian, &fat_cluster_count)
	if err != nil {
		fmt.Println("Error reading FAT cluster count:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("FAT cluster count read: %d\n", fat_cluster_count)

	// Read the total number of data clusters
	var cluster_count int32
	err = binary.Read(file, binary.LittleEndian, &cluster_count)
	if err != nil {
		fmt.Println("Error reading cluster count:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("Cluster count read: %d\n", cluster_count)

	// Read the starting positions
	var fat1_start int32
	err = binary.Read(file, binary.LittleEndian, &fat1_start)
	if err != nil {
		fmt.Println("Error reading FAT1 start position:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("FAT1 starts at: %d\n", fat1_start)

	var fat2_start int32
	err = binary.Read(file, binary.LittleEndian, &fat2_start)
	if err != nil {
		fmt.Println("Error reading FAT2 start position:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("FAT2 starts at: %d\n", fat2_start)

	var data_start int32
	err = binary.Read(file, binary.LittleEndian, &data_start)
	if err != nil {
		fmt.Println("Error reading data start position:", err)
		return 0, 0, 0, 0, 0, 0, 0
	}
	fmt.Printf("Data starts at: %d\n", data_start)

	fmt.Println("File system format loaded successfully!")
	return int(file_size), int(fat_size), int(fat_cluster_count), int(cluster_count), int(fat1_start), int(fat2_start), int(data_start)
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
	// TODO: move to the inti, to this point

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

	SaveFormatFile(filename, file_size, fat_size, fat_cluster_count, cluster_count, fat1_start, fat2_start, data_start)

	fmt.Println("File system formatted and saved successfully!")
	fmt.Println()
}
