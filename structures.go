package main

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
