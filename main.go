package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/ernestodeltoro/gofiledownload/progresswritter"
	ws "github.com/ernestodeltoro/gofiledownload/webscraper"
)

const (
	sourceFile int = iota
	appleFilex86
	appleFileARM
	linuxFile
	windowsFile
)

// FileData data struct to contain the info for downloading
type FileData struct {
	FilePath   string
	FileURL    string
	FileSHA256 string
}

// OsArch to return the type of platform the download program is running
type OsArch int

func main() {

	goVersion := GetInstalledGoVersion()
	fmt.Printf("Current go version: %s\n", goVersion)

	fd, err := InitializeFileData()
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err.Error())
		waitForEnterPress()
		return
	}

	fmt.Printf("To download:\n%s\n", fd.FileURL)

	start := time.Now()

	err = DownloadFile(fd)
	if err != nil {
		fmt.Println(err.Error())
		waitForEnterPress()
		return
	}

	elapsed := time.Since(start)

	fmt.Printf("Done in %s...\n", elapsed)

	err = VerifyFileSHA256(fd)
	if err != nil {
		fmt.Println(err.Error())
		waitForEnterPress()
		return
	}

	fmt.Println("SHA256 value verified, ok")
	waitForEnterPress()
}

func waitForEnterPress() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to exit: ")
	reader.ReadString('\n')
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(fd FileData) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(fd.FilePath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(fd.FileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return errors.New("Bad status: " + resp.Status)
	}

	pw := progresswritter.NewNonBloking(uint64(resp.ContentLength), 1000)

	// Write the body to file
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Closing the file before renaming it
	err = out.Close()
	if err != nil {
		return err
	}

	err = os.Rename(fd.FilePath+".tmp", fd.FilePath)
	if err != nil {
		return err
	}

	return nil
}

// VerifyFileSHA256 Compares the expected SHA256 with the actual value
// extracted from the file and returns true if they match
func VerifyFileSHA256(fd FileData) error {
	f, err := os.Open(fd.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return err
	}

	strSHA := fmt.Sprintf("%x", h.Sum(nil))

	if strSHA != fd.FileSHA256 {
		return errors.New("SHA values don't match")
	}

	return nil
}

// DownloadData will return the data needed to download and save the file
func DownloadData(osARCH OsArch) (fd FileData, err error) {

	homePage := "https://go.dev"
	seedURL := homePage + "/dl/"
	const numberOfHighlightedItemsToRetrieve = 5

	// Get the data
	resp, err := http.Get(seedURL)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer resp.Body.Close()

	links, err := ws.GetHighlightClassTokensN(resp, numberOfHighlightedItemsToRetrieve)
	if err != nil {
		return
	}

	if len(links) != numberOfHighlightedItemsToRetrieve {
		err = errors.New("unable to retrieve all the items from the download page")
		return
	}

	fd.FilePath = links[osARCH].FileName()
	fd.FileURL = makeProperHREF(links[osARCH].Href())
	fd.FileSHA256 = links[osARCH].Sha256()
	err = nil

	return
}

// makeProperHREF construct the proper HREF based on the home web page
func makeProperHREF(href string) string {
	// Make sure the url begins in http**
	hasProto := ws.HasHTTP(href)
	if hasProto {
		return href
	}

	newHref := ws.AddProto(href, "https://go.dev")
	return newHref
}

// GetOSFileIndex returns the detected OS type
func GetOSFileIndex() OsArch {
	var osARCH int
	switch runtime.GOOS {
	case "windows":
		osARCH = windowsFile
	case "linux":
		osARCH = linuxFile
	case "darwin":
		osARCH = appleFilex86
	default:
		osARCH = sourceFile
	}

	return OsArch(osARCH)
}

// GetInstalledGoVersion returns the curently installed go version
func GetInstalledGoVersion() string {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "Undefined / unable to determine"
	}
	return string(out)
}
