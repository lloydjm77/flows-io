package fileutil

import (
	"io"
	"net/http"
	"os"
)

// DetectFileType determines and returns the MIME type for the supplied file.
func DetectFileType(path string) string {
	file, _ := os.Open(path)

	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return err.Error()
	}

	// Reset the read pointer if necessary.
	file.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer[:n])

	return contentType
}
