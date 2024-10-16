package main

import (
	"errors"
	"fmt"
	"os"
)

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
func (fs *FileSystem) FindFreeCluster() (int, error) {

	for i, val := range fs.FatTable {
		if val == FAT_FREE {
			return i, nil
		}
	}

	return -1, errors.New("no free clusters available")
}

// Load file from external system into pseudo-FAT (incp s1 s2)
func (fs *FileSystem) InCp(filename string, data []byte) error {

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
func (fs *FileSystem) Cat(filename string) error {

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
func (fs *FileSystem) Rm(filename string) error {
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

func SaveFileSystem(filename string, fs *FileSystem) {

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write the FAT table
	for _, val := range fs.FatTable {
		_, err = file.Write([]byte{byte(val)})
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// Write the directory
	for _, entry := range fs.Directory {
		_, err = file.Write([]byte(entry.Name))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		_, err = file.Write([]byte{0})
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		_, err = file.Write([]byte{byte(entry.Size)})
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		_, err = file.Write([]byte{byte(entry.FirstCluster)})
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// Write the cluster data
	for _, data := range fs.ClusterData {
		_, err = file.Write(data)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	fmt.Println("File system saved successfully!")

}

func LoadFileSystem(filename string) *FileSystem {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	fs := &FileSystem{}
	fs.Init()

	// Read the FAT table
	for i := range fs.FatTable {
		var val byte
		_, err = file.Read([]byte{val})
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		fs.FatTable[i] = int(val)
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
		entry.Name = name
		_, err = file.Read([]byte{byte(entry.Size)})
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		_, err = file.Read([]byte{byte(entry.FirstCluster)})
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		fs.Directory[name] = entry
	}

	// Read the cluster data
	for i := range fs.ClusterData {
		_, err = file.Read(fs.ClusterData[i])
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
	}

	fmt.Println("File system loaded successfully!")
	return fs
}
