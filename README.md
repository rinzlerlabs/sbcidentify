# sbcidentity

A simple library for identifying the SBC on which the library is running. See the [cli](cmd/sbcidentify/main.go) for a simple example of usage. This can be used as an import or install the CLI version.

Currently supported boards:
* Raspberry Pis
* Various Jetson boards

## Package

Install the package with
```
go get github.com/rinzlerlabs/sbcidentify@latest
```

To identify a board, simply import
```
"github.com/rinzlerlabs/sbcidentify"
```
Then to identify the board
```
board, err := sbcidentify.GetBoardType()
```

## CLI

To install the CLI version, simply run
```
go install github.com/rinzlerlabs/sbcidentify/cmd/sbcidentify@latest
```

Usage
```
Usage of sbcidentify:
  -d    Enable debug logging
  -o string
        Specify the log output, accept StdOut, StdErr, or a file path (default "StdOut")
```
