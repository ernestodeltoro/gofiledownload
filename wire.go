//+build wireinject

package main

import "github.com/google/wire"

// InitializeEvent creates an Event
func InitializeFileData() (FileData, error) {
	wire.Build(DownloadData, GetOSFileIndex)
	return FileData{}, nil
}
