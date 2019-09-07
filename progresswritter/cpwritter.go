package progresswritter

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

// ConcurrentProgressWriter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type ConcurrentProgressWriter struct {
	currentWritten uint64
	fullSize       uint64
	newWrite       chan bool
}

// NewConcurrent returns a Writer interface that allows to show the
// writting progress in percentage given the fullSize of the file is known
func NewConcurrent(fullSize uint64) *ConcurrentProgressWriter {
	pw := ConcurrentProgressWriter{
		currentWritten: 0,
		fullSize:       fullSize,
		newWrite:       make(chan bool),
	}

	go pw.serveUpdateChannel()

	return &pw
}

func (pw *ConcurrentProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.currentWritten += uint64(n)
	pw.newWrite <- true
	return n, nil
}

func (pw *ConcurrentProgressWriter) serveUpdateChannel() {
	for range pw.newWrite {
		pw.updateProgress()
	}
}

// updateProgress prints
func (pw *ConcurrentProgressWriter) updateProgress() {

	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 40))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete of %s", humanize.Bytes(pw.currentWritten), humanize.Bytes(pw.fullSize))
}
