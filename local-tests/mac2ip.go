package main

import (
	"fmt"
	"leaf/tools"
	"time"
)

func main() {
	fmt.Println(tools.Mac2Ip("00:00:00:00:00:00"))
	fmt.Println(len("4aba0904a953_"))
	fmt.Println(len(time.RFC3339))
	fmt.Println(len("[07:57:11]"))

}
