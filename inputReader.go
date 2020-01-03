package main

import (
	"os"
	"bufio"
)

var lineScanner = bufio.NewScanner(os.Stdin)
var currentLine []byte
var currentCharIdx = 0

func readStdinChar() byte {
	if currentCharIdx >= len(currentLine) {
		lineScanner.Scan()
		currentLine = append(lineScanner.Bytes(), 10)
		currentCharIdx = 0
	}
	b := currentLine[currentCharIdx]
	currentCharIdx++
	return b
}