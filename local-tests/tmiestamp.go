package main

import (
	"fmt"
	"time"
)

func main() {

	StartNs := 7773564863511465
	fmt.Println(time.Unix(0, int64(StartNs)).Format("15:04:05"))

}
