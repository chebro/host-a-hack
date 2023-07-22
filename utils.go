package main

import (
	"fmt"
	"os"
)

func WriteStringToFile(filename, content string) error {
	// Open the file in write-only mode with file creation permission
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
