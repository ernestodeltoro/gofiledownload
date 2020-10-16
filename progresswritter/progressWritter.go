package progresswritter

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

// ProgressWriter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type ProgressWriter struct {
	currentWritten uint64
	fullSize       uint64
}

// New returns a Writer interface that allows to show the
// writing progress in percentage given the fullSize of the file is known
func New(fullSize uint64) *ProgressWriter {
	pw := ProgressWriter{0, fullSize}
	return &pw
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.currentWritten += uint64(n)
	pw.UpdateProgress()
	return n, nil
}

// UpdateProgress prints
func (pw *ProgressWriter) UpdateProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(pw.currentWritten))
}
