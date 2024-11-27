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

To check if a board is a specific type for hardware specific code, you can use `sbcidentify.IsBoardType()`. The boards definitions are structured such that they go from least to most restrictive.

For example, if you have code that should _only_ run on Raspberry Pi boards, you can do
```
if sbcidentify.IsBoardType(raspberrypi.RaspberryPi) {
	// do Pi specifc stuff
}
```

If you need to be more restrictive, say, only Raspberry Pi 5 boards
```
if sbcidentify.IsBoardType(raspberrypi.RaspberryPi5) {
	// do Pi 5 specifc stuff
} else if sbcidentify.IsBoardType(raspberrypi.RaspberryPi4) {
  // do Pi 4 specific stuff
} else {
  // fallback
}
```

The device tree heirarchies look like:
```
Raspberry Pi
├── 3
│   ├── 3B
│   │   └── 3BPlus
│   └── 3APlus
├── 4
│   └── 4B
│       ├── 4B 1GB
│       ├── 4B 2GB
│       ├── 4B 4GB
│       ├── 4B 8GB
│       └── 4400
└── 5
    └── 5B
        ├── 5B 2GB
        ├── 5B 4GB
        └── 5B 8GB

NVIDIA
├── Jetson
│   ├── Nano
│   │   ├── Nano 2GB
│   │   └── Nano 4GB
│   ├── TX1
│   ├── TX2
│   │   ├── TX2 4GB
│   │   ├── TX2 8GB
│   │   ├── TX2 NX
│   │   └── TX2i
│   ├── Xavier
│   │   ├── Xavier NX
│   │   └── AGX Xavier
│   └── Orin
│       ├── Orin NX
│       ├── Orin Nano
│       └── AGX Orin
├── Clara AGX
└── Shield TV
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
