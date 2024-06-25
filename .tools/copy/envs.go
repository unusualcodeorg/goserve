package main

import (
	"fmt"
	"io"
	"os"
)

func copyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}

	// Ensure all data is flushed to the destination file
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("could not flush data to destination file: %w", err)
	}

	return nil
}

func main() {
	err := copyFile(".env.example", ".env")
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}

	err = copyFile(".test.env.example", ".test.env")
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}

	fmt.Println("env files copied successfully!")
}
