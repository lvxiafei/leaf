package main

import (
	"fmt"
	"strings"
)

func main() {
	inputString := "Hello, World! Hello, Go!"

	// 将 "Hello" 替换为 "Hi"
	newString := strings.Replace(inputString, "Hello", "Hi", -1)
	fmt.Println("Modified String:", newString)

	release := "4.18.0-425.3.1.el8.x86_64"
	newString = strings.Replace(release, ".x86_64", "", -1)

	fmt.Println("Modified release:", newString)
}
