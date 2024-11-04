package main

// Constants for file system
const (
	CLUSTER_SIZE  = 1024 // 1KB cluster size
	FAT_ENTRY     = 4    // FAT entry size in bytes
	MAX_FILE_NAME = 12   // 8.3 format = 11 chars + null terminator
	FAT_FREE      = -1   // FAT free cluster marker
	FAT_EOF       = -2   // FAT end of file marker
)

// FileSystemFormat struct to store file system metadata
type FileSystemFormat struct {
	file_size         int32
	fat_size          int32
	fat_cluster_count int32
	cluster_count     int32
	fat1_start        int32
	fat2_start        int32
	data_start        int32
}

// FAT entry struct to simulate FAT table
type FAT []int

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
