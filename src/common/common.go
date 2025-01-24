package common

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

var (
	errUuidToBinaryInvalid       = errors.New("invalid uuid format")
	errBinaryToUuidInvalidBytes  = errors.New("invalid bytes")
	errBinaryToUuidCouldNotParse = errors.New("could not parse bytes to uuid")
)

func UuidToBinary(uuid uuid.UUID) ([]byte, error) {
	id := uuid.String()
	hexStr := id[0:8] + id[9:13] + id[14:18] + id[19:23] + id[24:]
	binary, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, errUuidToBinaryInvalid
	}
	return binary, nil
}

func BinaryToUuid(bytes []byte) (uuid.UUID, error) {
	var binary []byte
	nbOfBytes := hex.Encode(binary, bytes)
	if nbOfBytes == 0 {
		return uuid.Nil, errBinaryToUuidInvalidBytes
	}

	id, err := uuid.ParseBytes(binary)

	if err != nil {
		return uuid.Nil, errBinaryToUuidCouldNotParse
	}

	return id, nil
}

func ReadFileText(fileName string) (string, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return "", fmt.Errorf("could not open file %s", fileName)
	}

	defer file.Close()
	var content string

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file %s", fileName)
	}

	return content, nil
}

func PrintFormat(msg string, a ...any) {
	fmt.Println(fmt.Sprintf(msg, a...))
}

func RequestBodyToString(reader io.ReadCloser) (string, error) {
	buff := new(strings.Builder)
	if _, err := io.Copy(buff, reader); err != nil {
		return "", err
	}

	return buff.String(), nil
}
