## goFileDownload
Just code for downloading Go compiler files with Golang

[![Go Report Card](https://goreportcard.com/badge/github.com/ernestodeltoro/goFileDownload)](https://goreportcard.com/report/github.com/ernestodeltoro/gofiledownload)


Very simple program that search for the "install" file on [https://golang.org/dl/](https://golang.org/dl/) based on your O.S. Makes the download, and verifies it's checksum.

## Dependencies and Acknowlegments

I'm using google/wire: [https://github.com/google/wire](https://github.com/google/wire), which is a Compile-time Dependency Injection for Go. I know this is an overkill, but I just wanted to tested it. In any case this could also serve as a wire example. 

Probably you will need to install wire in order to build this code:

```shell
$ go get github.com/google/wire/cmd/wire
```

and ensuring that `$GOPATH/bin` is added to your `$PATH`.

I'm also using some code extracted from [https://golangcode.com/download-a-file-with-progress/](https://golangcode.com/download-a-file-with-progress/)

## Running it

Download the source code:
```shell
$ git clone https://github.com/ernestodeltoro/goFileDownload.git
```

Go to the download folder and run it:
```shell
$ go run main.go
```
The installation file will download to the same folder.

## Contribution

Thank you for considering to help out with the source code! I am grateful for even the smallest of fixes or improvements! I'm lacking tests for this application, so if you want you can start there.

If you'd like to contribute to goFileDownload, follow these steps:

1. Fork repository ([github help](https://help.github.com/en/articles/fork-a-repo))
2. Create a branch
3. Make the changes
4. Commit and send a pull request

I will review it and merge into the main code base.