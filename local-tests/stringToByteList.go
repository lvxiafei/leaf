package main

import (
	"fmt"
	"leaf/tools"
)

func main() {
	// 输入字符串
	str := "0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 90 70 95 0 159 255 255 19 244 7"

	// 调用封装的函数
	byteSlice, err := tools.StringToByteArray(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 输出 byte 切片
	fmt.Println(byteSlice)
}
