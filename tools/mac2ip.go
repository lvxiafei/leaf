package tools

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func Mac2Ip(macAddress string) string {
	ip, err := findIPByMAC(macAddress)
	if err != nil {
		ipLocal, err := IfaceMacToIP(macAddress)
		if err != nil {
			return "0.0.0.0"
		} else {
			return ipLocal
		}

	} else {
		return ip
	}
}

func findIPByMAC(macAddress string) (string, error) {
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 6 && fields[3] == macAddress {
			return fields[0], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("MAC address not found in ARP table")
}

func IfaceMacToIP(mac string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		// 获取接口的硬件地址（MAC 地址）
		hwAddr := iface.HardwareAddr.String()
		//fmt.Println(hwAddr)
		if hwAddr == mac {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			// 通常情况下，一个接口可能有多个 IP 地址，这里只返回第一个
			if len(addrs) > 0 {
				ip, _, err := net.ParseCIDR(addrs[0].String())
				if err != nil {
					return "", err
				}
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("MAC address not found in iface")
}

func MacToUpperCaseString(mac [6]uint8) string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func MacToLowerCaseString(mac [6]uint8) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}
