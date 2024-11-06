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

	err := binary.Read(file, binary.LittleEndian, value)
	if err != nil {
		fmt.Println("Error reading from file:", err)
	}
}
