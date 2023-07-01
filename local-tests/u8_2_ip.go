package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	// 示例 daddr 地址
	daddr := uint64(0x7f9d5ef22ee4)

	// 使用 net.IP 类型直接创建 IP 地址
	ip := make(net.IP, 4)
	binary.LittleEndian.PutUint32(ip, uint32(daddr))

	// 打印 IP 地址
	fmt.Println("IPv4 Address:", ip)
}
