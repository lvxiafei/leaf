package tools

import (
	"fmt"
	"strconv"
	"strings"
)

func StringToByteArray(str string) ([]byte, error) {
	// 使用 strings.Fields 分割字符串
	strList := strings.Fields(str)

	// 转换字符串切片为 byte 切片
	var byteSlice []byte
	for _, s := range strList {
		// 将字符串转换为整数
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("Error converting string to integer: %v", err)
		}
		// 将整数转换为 byte，并添加到切片
		byteSlice = append(byteSlice, byte(num))
	}

	return byteSlice, nil
}
