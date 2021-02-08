package progresswritter

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// NonBlokingProgressWriter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type NonBlokingProgressWriter struct {
	currentWritten uint64
	fullSize       uint64
	newWrite       chan uint64
	sleepTime      time.Duration
}

// NewNonBloking returns a Writer interface that allows to show the
// writing progress in percentage given the fullSize of the file is known
func NewNonBloking(fullSize uint64, sleepTime time.Duration) *NonBlokingProgressWriter {
	pw := NonBlokingProgressWriter{
		currentWritten: 0,
		fullSize:       fullSize,
		newWrite:       make(chan uint64),
		sleepTime:      sleepTime,
	}

	go pw.serveUpdateChannel()

	return &pw
}

func (pw *NonBlokingProgressWriter) serveUpdateChannel() {

	defer close(pw.newWrite)

	for nw := range pw.newWrite {
		pw.updateProgress(nw)
		if nw >= pw.fullSize {
			return
		}
		time.Sleep(pw.sleepTime * time.Millisecond)
	}
}

// Non Blocking write on a channel
func (pw *NonBlokingProgressWriter) execUpdateNonBlocking(currentWritten uint64) {
	select {
	case pw.newWrite <- currentWritten:
		return
	default:
	}
}

func (pw *NonBlokingProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.currentWritten += uint64(n)
	pw.execUpdateNonBlocking(pw.currentWritten)
	return n, nil
}

// updateProgress prints
func (pw *NonBlokingProgressWriter) updateProgress(currentWritten uint64) {

	fmt.Printf("\r%s", strings.Repeat(" ", 40))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete of %s", humanize.Bytes(currentWritten), humanize.Bytes(pw.fullSize))

}
