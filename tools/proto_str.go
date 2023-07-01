package tools

import (
	"encoding/binary"
	"unsafe"
)

/*
func EthProtoString(proto int) string {
	// ether proto definitions:
	// https://sites.uclouvain.be/SystInfo/usr/include/linux/if_ether.h.html
	// IEEE 802 Numbers https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml
	// https://pkg.go.dev/syscall
	protoStr := fmt.Sprintf("UNKNOWN#%d", proto)
	switch proto {
	case syscall.ETH_P_ALL:
		protoStr = "ALL"
	case syscall.ETH_P_IP: // Ox0800,2048
		protoStr = "IP"
	case syscall.ETH_P_ARP: // 0x0806,2054
		protoStr = "ARP"
	case syscall.ETH_P_RARP:
		protoStr = "RARP"
	case syscall.ETH_P_IPV6:
		protoStr = "IPV6"
	case syscall.ETH_P_8021Q: // 0x8100,33024
		protoStr = "VLAN"
	}
	return protoStr
}
*/

// Htons converts the unsigned short integer hostshort from host byte order to network byte order.
func Htons(i uint16) uint16 {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return *(*uint16)(unsafe.Pointer(&b[0]))
}
