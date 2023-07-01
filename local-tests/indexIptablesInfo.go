package main

import (
	"bytes"
	"encoding/binary"
	"leaf/tools"
	"log"
	"unsafe"
)

type l2Info struct {
	DestMac [6]byte
	L3Proto uint16
}

type l3Info struct {
	Saddr     [16]byte
	Daddr     [16]byte
	TotLen    uint16
	IPVersion uint8
	L4Proto   uint8
}

type l4Info struct {
	Sport    uint16
	Dport    uint16
	TCPFlags uint16
	Pad      [2]byte
}

type icmpInfo struct {
	IcmpID   uint16
	IcmpSeq  uint16
	IcmpType uint8
	Pad      [3]byte
}

type iptablesInfo struct {
	TableName [32]byte
	Verdict   uint32
	IptDelay  uint64
	Hook      uint8
	Pf        uint8
	Pad       [2]byte
}

type iptablesTrace struct {
	In        [16]byte
	Out       [16]byte
	TableName [32]byte
	ChainName [32]byte
	RuleNum   uint32
	HookNum   uint32
	Pf        uint8
	Pad       [3]uint8
}

type pktInfo struct {
	Ifname  [16]byte
	Len     uint32
	CPU     uint32
	Pid     uint32
	NetNS   uint32
	PktType uint8
	Pad     [3]byte
}

type perfEvent struct {
	// order
	FuncName [32]byte

	Skb     uint64
	StartNs uint64
	Flags   uint8
	Pad     [3]byte

	pktInfo
	l2Info
	l3Info
	l4Info
	icmpInfo
	iptablesInfo
}

func main() {
	var event perfEvent
	var ipt iptablesInfo
	//indexIptablesInfo := unsafe.Sizeof(perfEvent{})
	indexIptablesInfo := unsafe.Sizeof(perfEvent{}) - unsafe.Sizeof(iptablesInfo{}) // 152
	log.Printf("indexIptablesInfo: %+v", indexIptablesInfo)

	str := "0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 90 70 95 0 159 255 255 19 244 74 245 186 26 6 0 3 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 84 0 0 0 0 0 0 0 142 202 0 0 152 0 0 240 0 0 0 0 0 0 0 0 0 0 0 8 10 142 53 176 0 0 0 0 0 0 0 0 0 0 0 0 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 84 0 4 1 0 0 0 0 0 0 0 0 87 0 1 0 8 0 0 0 114 97 119 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 127 113 0 0 0 0 0 0 3 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0"
	RawSample, _ := tools.StringToByteArray(str)
	if err := binary.Read(bytes.NewBuffer(RawSample), binary.LittleEndian, &event); err != nil {
		log.Printf("Failed to parse perf event: %v", err)
	}

	log.Printf("record.RawSample: %+v", RawSample)
	log.Printf("record.RawSample indexIptablesInfo: %+v", RawSample[152:])
	log.Printf("event.iptablesInfo: %+v", event.iptablesInfo)

	if err := binary.Read(bytes.NewReader(RawSample[indexIptablesInfo-4:]), binary.LittleEndian, &ipt); err != nil {
		log.Printf("Failed to parse iptables info: %v", err)
	}

	log.Printf("iptablesInfo: %+v", ipt)

}
