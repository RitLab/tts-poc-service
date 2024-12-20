package domain

import (
	"fmt"
	"io"
	"os"
)

// AppendFile reads data from an input file and writes it to the output writer
func AppendFile(out io.Writer, file string) error {
	// Open the input file
	in, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer in.Close()

	// Copy the file data to the output
	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	return nil
}
