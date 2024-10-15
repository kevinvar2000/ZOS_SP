package main

import (
	"errors"
	"fmt"
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
