package main

// Constants for file system
const (
	CLUSTER_SIZE  = 1024 // 1KB cluster size
	FAT_ENTRY     = 4    // FAT entry size in bytes
	MAX_FILE_NAME = 12   // 8.3 format = 11 chars + null terminator
	FAT_FREE      = -1   // FAT free cluster marker
	FAT_EOF       = -2   // FAT end of file marker
	FAT_BAD       = -3   // FAT bad cluster marker
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
	Name          [MAX_FILE_NAME]byte
	Size          int32
	First_cluster int32
	Is_directory  uint8 // use 1 for true and 0 for false
}
