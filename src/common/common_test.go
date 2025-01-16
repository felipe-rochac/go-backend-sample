package common

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestReadFileText(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some content to the temporary file
	content := "Hello, World!\nThis is a test file."
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test ReadFileText function
	readContent, err := ReadFileText(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadFileText returned an error: %v", err)
	}

	if readContent != content+"\n" {
		t.Errorf("ReadFileText returned unexpected content: got %v, want %v", readContent, content)
	}

	// Test ReadFileText with a non-existent file
	_, err = ReadFileText("non_existent_file.txt")
	if err == nil {
		t.Error("Expected an error for non-existent file, but got nil")
	}
}
func Test_UuidToBinary_ReturnBinary(t *testing.T) {
	uuidStr := "123e4567-e89b-12d3-a456-426614174000"
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Fatalf("Failed to parse UUID: %v", err)
	}

	binary, err := UuidToBinary(u)
	if err != nil {
		t.Fatalf("UuidToBinary returned an error: %v", err)
	}

	expectedBinary, _ := hex.DecodeString("123e4567e89b12d3a456426614174000")
	if !bytes.Equal(binary, expectedBinary) {
		t.Errorf("UuidToBinary returned unexpected binary: got %v, want %v", binary, expectedBinary)
	}
}

func Test_BinaryToUuid_ExpectUuid(t *testing.T) {
	binary, _ := hex.DecodeString("123e4567e89b12d3a456426614174000")
	u, err := BinaryToUuid(binary)
	if err != nil {
		t.Fatalf("BinaryToUuid returned an error: %v", err)
	}

	expectedUuid := "123e4567-e89b-12d3-a456-426614174000"
	if u.String() != expectedUuid {
		t.Errorf("BinaryToUuid returned unexpected UUID: got %v, want %v", u.String(), expectedUuid)
	}

	// Test with invalid bytes
	_, err = BinaryToUuid([]byte{0x00})
	if err == nil {
		t.Error("Expected an error for invalid bytes, but got nil")
	}
}

func Test_BinaryToUuid_InvalidUuid_ExpectError(t *testing.T) {
	binary, _ := hex.DecodeString("123e4567e89b12d3a4564")
	_, err := BinaryToUuid(binary)
	if err == nil {
		t.Fatalf("BinaryToUuid returned an error: %v", err)
	}
}

func Test_PrintFormat_ExpectSuccess(t *testing.T) {
	// Capture the output of PrintFormat
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintFormat("Hello, %s!", "World")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	expectedOutput := "Hello, World!\n"
	if buf.String() != expectedOutput {
		t.Errorf("PrintFormat returned unexpected output: got %v, want %v", buf.String(), expectedOutput)
	}
}

func Test_BinaryToUuid_EmptySlice_ExpectError(t *testing.T) {
	// Test with empty byte slice
	_, err := BinaryToUuid([]byte{})
	if err == nil {
		t.Error("Expected an error for empty byte slice, but got nil")
	} else if err != errBinaryToUuidInvalidBytes {
		t.Errorf("Expected error %v, but got %v", errBinaryToUuidInvalidBytes, err)
	}

	// Test with invalid byte slice length
	_, err = BinaryToUuid([]byte{0x12, 0x34})
	if err == nil {
		t.Error("Expected an error for invalid byte slice length, but got nil")
	} else if err != errBinaryToUuidInvalidBytes {
		t.Errorf("Expected error %v, but got %v", errBinaryToUuidInvalidBytes, err)
	}

	// Test with invalid byte content
	_, err = BinaryToUuid([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err == nil {
		t.Error("Expected an error for invalid byte content, but got nil")
	} else if err != errBinaryToUuidCouldNotParse {
		t.Errorf("Expected error %v, but got %v", errBinaryToUuidCouldNotParse, err)
	}
}
