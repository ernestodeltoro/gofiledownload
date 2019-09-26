package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/ernestodeltoro/goFileDownload/progresswritter"
)

func main() {

	filePath, fileURL, fileSHA256, err := DownloadData()
	if err != nil {
		fmt.Printf(err.Error())
	}

	goVersion := runtime.Version()
	fmt.Printf("Current go version %s\n", goVersion)
	fmt.Printf("Downloading: %s\n", filePath)

	err = DownloadFile(filePath, fileURL)
	if err != nil {
		fmt.Printf(err.Error())
	}

	shaOK, err := VerifyFileSHA256(filePath, fileSHA256)
	if err != nil {
		fmt.Printf(err.Error())
	}

	if !shaOK {
		fmt.Println("SHA256 values don't match")
	} else {
		fmt.Println("SHA256 value verified")
	}

}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// extracted from: https://golangcode.com/download-a-file-with-progress/
func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return errors.New("Bad status: " + resp.Status)
	}

	pw := progresswritter.NewNonBloking(uint64(resp.ContentLength))

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

	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}

// VerifyFileSHA256 Compares the expected SHA256 with the actual value
// extracted from the file and returns true if they match
func VerifyFileSHA256(filePath, expectedFileSHA256 string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return false, err
	}

	strSHA := fmt.Sprintf("%x", h.Sum(nil))

	if strSHA != expectedFileSHA256 {
		return false, nil
	}

	return true, nil
}

// DownloadData will return the data needed to download and save the file
func DownloadData() (filePath, fileURL, fileSHA256 string, err error) {
	filePath = "go1.13.1.windows-amd64.msi"
	fileURL = "https://dl.google.com/go/go1.13.1.windows-amd64.msi"
	fileSHA256 = "ee6ff4f31766096e6ff00bdd855e007c473cc9080f3588d7be0c5ec0e0b85578"
	err = nil

	return
}
