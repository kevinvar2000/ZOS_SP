package main

// Constants for file system
const (
	CLUSTER_SIZE      = 1024 // 1KB cluster size
	MAX_FILE_NAME     = 12   // 8.3 format = 11 chars + null terminator
	MAX_CLUSTER_COUNT = 1024 // Max number of clusters in the file system
	FAT_FREE          = -1   // FAT free cluster marker
	FAT_EOF           = -2   // FAT end of file marker
	MAX_FILE_LEN      = 5000 // Maximum file length for the `short` command
	FAT_ENTRY         = 4    // FAT entry size in bytes
)

// FAT entry struct to simulate FAT table
type FAT [MAX_CLUSTER_COUNT]int

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

type FileInfo struct {
	Name string
	Size int
}
