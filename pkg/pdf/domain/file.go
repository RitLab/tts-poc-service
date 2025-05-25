package domain

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	pkgError "tts-poc-service/pkg/common/error"
)

var PdfExtension = "pdf"
var EmptyFile = fmt.Errorf(pkgError.FILE_IS_EMPTY)
var WrongFileExtension = fmt.Errorf(pkgError.FILE_EXTENSION_NOT_SUPPORTED)

// ValidatePdfFile reads file format based on file name extension
func ValidatePdfFile(file *multipart.FileHeader) error {
	if file == nil {
		return EmptyFile
	}
	ext := strings.Split(filepath.Ext(file.Filename), ".")
	if len(ext) < 2 {
		return WrongFileExtension
	}

	if ext[1] != PdfExtension {
		return WrongFileExtension
	}
	return nil
}
