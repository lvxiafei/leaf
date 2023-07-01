package main

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)

func examples() {

	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println(dmp.DiffCleanupMerge(diffs))
}

func main() {
	examples()
}
