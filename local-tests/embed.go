package main

import (
	"embed"
	"fmt"
)

//go:embed *
var f embed.FS

func funcDir1() {
	dirEntries, _ := f.ReadDir("draw_json")
	for _, de := range dirEntries {
		fmt.Println(de.Name(), de.IsDir())
	}
}

//go:embed draw_json
var f2 embed.FS

func funcDir2() {

	data, _ := f2.ReadFile("draw_json/copy_text.json")
	fmt.Println(string(data))
	data, _ = f2.ReadFile("draw_json/main.go")
	fmt.Println(string(data))
}

func main() {
	//funcDir1()
	funcDir2()
}
