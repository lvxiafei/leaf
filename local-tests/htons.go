package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

func Htons(i uint16) uint16 {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return *(*uint16)(unsafe.Pointer(&b[0]))
}
func main() {
	ret := int(Htons(4000))
	fmt.Println(ret)
}
