package domain

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	pkgError "tts-poc-service/pkg/common/error"
)

var AudioExtension = "mp3"
var EmptyFile = fmt.Errorf(pkgError.FILE_IS_EMPTY)
var WrongFileExtension = fmt.Errorf(pkgError.FILE_EXTENSION_NOT_SUPPORTED)

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

// ValidateAudioFile reads file format based on file name extension
func ValidateAudioFile(file *multipart.FileHeader) error {
	if file == nil {
		return EmptyFile
	}
	ext := strings.Split(filepath.Ext(file.Filename), ".")
	if len(ext) < 2 {
		return WrongFileExtension
	}

	if ext[1] != AudioExtension {
		return WrongFileExtension
	}
	return nil
}
