package main

// Constants for file system
const (
	CLUSTER_SIZE      = 1024 // 1KB cluster size
	FAT_ENTRY         = 4    // FAT entry size in bytes
	MAX_FILE_NAME     = 12   // 8.3 format = 11 chars + null terminator
	MAX_CLUSTER_COUNT = 1024 // Max number of clusters in the file system
	FAT_FREE          = 0    // FAT free cluster marker
	FAT_EOF           = -1   // FAT end of file marker
)

// FAT entry struct to simulate FAT table
type FAT [MAX_CLUSTER_COUNT]int

// DirectoryEntry stores file metadata
type DirectoryEntry struct {
	name          string
	size          int
	first_cluster int
	is_directory  bool
}

// FileSystem struct
type FileSystem struct {
	fat_table    FAT
	directory    map[string]DirectoryEntry
	cluster_data [][]byte // Storage for the actual file data
}

type FileInfo struct {
	name string
	size int
}
