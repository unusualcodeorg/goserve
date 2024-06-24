package utils

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestLoadPEMFileInto(t *testing.T) {
    // Create a temporary PEM file for testing
    pemData := []byte("mock PEM file content")
    tmpfile, err := os.CreateTemp("", "example")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name()) // Clean up

    // Write PEM data to the temporary file
    if _, err := tmpfile.Write(pemData); err != nil {
        tmpfile.Close()
        t.Fatal(err)
    }
    if err := tmpfile.Close(); err != nil {
        t.Fatal(err)
    }

    // Test loading the PEM file into memory
    loadedData, err := LoadPEMFileInto(tmpfile.Name())
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Assert that the loaded data matches the original PEM data
    assert.Equal(t, pemData, loadedData)
}

func TestLoadPEMFileIntoError(t *testing.T) {
    // Test case where the file does not exist
    _, err := LoadPEMFileInto("nonexistent-file.pem")
    assert.Error(t, err)
}