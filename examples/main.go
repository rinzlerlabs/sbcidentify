package main

import (
	"fmt"

	"github.com/thegreatco/sbcidentify"
)

func main() {
	board, err := sbcidentify.GetBoardType()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Board type: %s\n", board)
}
