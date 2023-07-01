package main

import (
	"fmt"
	"strings"
)

func main() {
	if find := strings.Contains("test-v1", "v1"); find {
		fmt.Println("find the character.")
	}
}
