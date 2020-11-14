package progresswritter

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
)

// ConcurrentProgressWriter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type ConcurrentProgressWriter struct {
	currentWritten uint64
	fullSize       uint64
	fullSizeStr    string
	sleepTime      time.Duration
}

// NewConcurrent returns a Writer interface that allows to show the
// writing progress in percentage given the fullSize of the file is known
// fullSize is the size of the content to be downloaded, sleepTime is the
// time between os.Stdout updates
func NewConcurrent(fullSize uint64, sleepTime time.Duration) *ConcurrentProgressWriter {
	pw := ConcurrentProgressWriter{
		currentWritten: 0,
		fullSize:       fullSize,
		fullSizeStr:    humanize.Bytes(fullSize),
		sleepTime:      sleepTime,
	}

	go pw.serveUpdateChannel()

	return &pw
}

func (pw *ConcurrentProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.currentWritten += uint64(n)
	if atomic.CompareAndSwapUint64(&pw.currentWritten, pw.fullSize, pw.fullSize) {
		// print for the last time
		pw.updateProgress()
	}
	return n, nil
}

func (pw *ConcurrentProgressWriter) serveUpdateChannel() {
	for {
		time.Sleep(pw.sleepTime * time.Millisecond)
		if atomic.CompareAndSwapUint64(&pw.currentWritten, pw.fullSize, pw.fullSize) {
			return
		}
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
	fmt.Printf("\rDownloading... %s complete of %s", humanize.Bytes(pw.currentWritten), pw.fullSizeStr)
}
