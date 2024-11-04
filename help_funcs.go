package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func WriteToFile(file *os.File, value int32) {

	err := binary.Write(file, binary.LittleEndian, value)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

}

func ReadFromFile(file *os.File, value *int32) {

	err := binary.Read(file, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("Error reading from file:", err)
	}
}

// Helper function to get minimum value
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func PrintFileSystem(fs *FileSystem, outputFilename string) {

	file, err := os.Create(outputFilename)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString("File System Details:\n")
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(file, "Total Clusters: %d\n", len(fs.fat_table))
	if err != nil {
		return
	}

	_, err = file.WriteString("FAT Table:\n")
	if err != nil {
		return
	}
	for i, val := range fs.fat_table {
		_, err = fmt.Fprintf(file, "  Cluster %d: %d\n", i, val)
		if err != nil {
			return
		}
	}

	_, err = file.WriteString("Directory Entries:\n")
	if err != nil {
		return
	}
	for name, entry := range fs.directory {
		_, err = fmt.Fprintf(file, "  Name: %s, Size: %d bytes, First Cluster: %d\n", name, entry.size, entry.first_cluster)
		if err != nil {
			return
		}
	}

	_, err = file.WriteString("Cluster Data (first few bytes):\n")
	if err != nil {
		return
	}
	for i, data := range fs.cluster_data {
		_, err = fmt.Fprintf(file, "  Cluster %d: %x\n", i, data[:Min(len(data), 16)]) // Print first 16 bytes of each cluster
		if err != nil {
			return
		}
	}

}
