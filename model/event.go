package model

import (
	"fmt"
	"leaf/config"
	"leaf/tools"
	"net"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	ethProtoIP   = 0x0800
	ethProtoIPv6 = 0x86DD
)

const (
	IpprotoICMP   = 1
	IpprotoTCP    = 6
	IpprotoUDP    = 17
	ipprotoICMPv6 = 58
)

const (
	routeEventIf      = 0x0001
	RouteEventIptable = 0x0002

	RouteEventIptablesTrace = 0x0004
)

const (
	// the services chain
	kubeServicesChain string = "KUBE-SERVICES"

	// the external services chain
	kubeExternalServicesChain string = "KUBE-EXTERNAL-SERVICES"

	// the nodeports chain
	kubeNodePortsChain string = "KUBE-NODEPORTS"

	// the kubernetes postrouting chain
	kubePostroutingChain string = "KUBE-POSTROUTING"

	// KubeMarkMasqChain is the mark-for-masquerade chain
	KubeMarkMasqChain string = "KUBE-MARK-MASQ"

	// KubeMarkDropChain is the mark-for-drop chain
	KubeMarkDropChain string = "KUBE-MARK-DROP"

	// the kubernetes forward chain
	kubeForwardChain string = "KUBE-FORWARD"

	// kube proxy canary chain is used for monitoring rule reload
	kubeProxyCanaryChain string = "KUBE-PROXY-CANARY"
)

var (
	nfVerdictName = []string{
		"DROP",
		"ACCEPT",
		"STOLEN",
		"QUEUE",
		"REPEAT",
		"STOP",
	}

	hookNames = []string{
		"PREROUTING",
		"INPUT",
		"FORWARD",
		"OUTPUT",
		"POSTROUTING",

		// kube-proxy
		kubeServicesChain,
		kubeExternalServicesChain,
		kubeNodePortsChain,
		kubePostroutingChain,
		KubeMarkMasqChain,
		KubeMarkDropChain,
		kubeForwardChain,
		kubeProxyCanaryChain,
	}

	tcpFlagNames = []string{
		"CWR",
		"ECE",
		"URG",
		"ACK",
		"PSH",
		"RST",
		"SYN",
		"FIN",
	}
)

func _get(names []string, idx uint32, defaultVal string) string {
	if int(idx) < len(names) {
		return names[idx]
	}

	return defaultVal
}

type l2Info struct {
	DestMac [6]byte
	SrcMac  [6]byte

	EthProto uint32
	L3Proto  uint16

	VlanId uint16

	VlanTci   uint16
	VlanProto uint16
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

type IptablesInfo struct {
	TableName [32]byte
	Verdict   uint32
	IptDelay  uint64
	Hook      uint8
	Pf        uint8
	Pad       [2]byte
}

type IptablesTrace struct {
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
	Comm    [16]byte
}

type PerfEvent struct {
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
	IptablesInfo
}

type perfEventWithIptablesTrace struct {
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
	IptablesTrace
}

const (
	sizeofEvent = 116 // Note: 116 instead of int(unsafe.Sizeof(PerfEvent{})), because of alignment
	//sizeofEvent = 208
)

var earliestTs = uint64(0)

func (e *PerfEvent) outputTimestamp() string {
	if config.Cfg.Timestamp {
		if earliestTs == 0 {
			earliestTs = e.StartNs
		}
		return fmt.Sprintf("%-7.6f", float64(e.StartNs-earliestTs)/1000000000.0)
	}
	return time.Unix(0, int64(e.StartNs)).Format("15:04:05")
}

func (e *PerfEvent) outputTcpFlags() string {
	var flags []string
	tcpFlags := uint8(e.TCPFlags >> 8)
	for i := 0; i < len(tcpFlagNames); i++ {
		if tcpFlags&(1<<i) != 0 {
			flags = append(flags, tcpFlagNames[i])
		}
	}

	return strings.Join(flags, ",")
}

func (e *PerfEvent) outputIptablesInfo(ipt *IptablesInfo, trace *IptablesTrace) string {
	var sb strings.Builder

	pktType := e.outputPktType(e.PktType)

	if e.Flags&RouteEventIptable == RouteEventIptable {
		pf := "PF_INET"
		if ipt.Pf == 10 {
			pf = "PF_INET6"
		}

		iptName := nullTerminatedStr(ipt.TableName[:])
		hook := _get(hookNames, uint32(ipt.Hook), fmt.Sprintf("~UNK~[%d]", ipt.Hook))
		verdict := _get(nfVerdictName, ipt.Verdict, fmt.Sprintf("~UNK~[%d]", ipt.Verdict))
		cost := time.Duration(ipt.IptDelay)

		fmt.Fprintf(&sb, "pkt_type=%-9s iptables=[pf=%s table=%s hook=%s verdict=%s cost=%s]", pktType,
			pf, iptName, hook, verdict, cost)
	}
	if e.Flags&RouteEventIptablesTrace == RouteEventIptablesTrace {
		pf := "PF_INET"
		if trace.Pf == 10 {
			pf = "PF_INET6"
		}

		in, out := nullTerminatedStr(trace.In[:]), nullTerminatedStr(trace.Out[:])
		table, chain := nullTerminatedStr(trace.TableName[:]), nullTerminatedStr(trace.ChainName[:])

		fmt.Fprintf(&sb, "pkt_type=%-9s ipttrace=[pf=%s in=%s out=%s table=%s chain=%s hook=%d rulenum=%d]", pktType,
			pf, in, out, table, chain, trace.HookNum, trace.RuleNum)
	}
	if sb.String() == "" {
		funcName := NullTerminatedStringToString(e.FuncName[:])
		return fmt.Sprintf("pkt_type=%-9s func=%s", pktType, funcName)
	}

	return sb.String()
}

func (e *PerfEvent) outputPktInfo() string {
	var saddr, daddr net.IP
	if e.l2Info.L3Proto == ethProtoIP {
		saddr = net.IP(e.Saddr[:4])
		daddr = net.IP(e.Daddr[:4])
	} else {
		saddr = net.IP(e.Saddr[:])
		daddr = net.IP(e.Daddr[:])
	}

	if e.L4Proto == IpprotoTCP {
		tcpFlags := e.outputTcpFlags()
		if tcpFlags == "" {
			return fmt.Sprintf("T:%s:%d->%s:%d",
				saddr, e.Sport, daddr, e.Dport)
		}
		return fmt.Sprintf("T_%s:%s:%d->%s:%d", tcpFlags,
			saddr, e.Sport, daddr, e.Dport)

	} else if e.L4Proto == IpprotoUDP {
		return fmt.Sprintf("U:%s:%d->%s:%d",
			saddr, e.Sport, daddr, e.Dport)
	} else if e.L4Proto == IpprotoICMP || e.L4Proto == ipprotoICMPv6 {
		if e.IcmpType == 8 || e.IcmpType == 128 {
			return fmt.Sprintf("I_request:%s->%s", saddr, daddr)
		} else if e.IcmpType == 0 || e.IcmpType == 129 {
			return fmt.Sprintf("I_reply:%s->%s", saddr, daddr)
		} else {
			return fmt.Sprintf("I:%s->%s", saddr, daddr)
		}
	} else {
		return fmt.Sprintf("%d:%s->%s", e.L4Proto, saddr, daddr)
	}
}

func nullTerminatedStr(b []byte) string {
	off := 0
	for ; off < len(b) && b[off] != 0; off++ {
	}
	b = b[:off]
	return *(*string)(unsafe.Pointer(&b))
}

// NullTerminatedStringToString is helper to convert null terminated string to GO string
func NullTerminatedStringToString(val []byte) string {
	// Calculate null terminated string len
	slen := len(val)
	for idx, ch := range val {
		if ch == 0 {
			slen = idx
			break
		}
	}
	return string(val[:slen])
}

func (e *PerfEvent) outputTraceInfo() string {
	iptables := ""
	if e.Flags&RouteEventIptable == RouteEventIptable {
		pf := "PF_INET"
		if e.Pf == 10 {
			pf = "PF_INET6"
		}
		iptName := nullTerminatedStr(e.TableName[:])
		hook := _get(hookNames, uint32(e.Hook), fmt.Sprintf("~UNK~[%d]", e.Hook))
		verdict := _get(nfVerdictName, e.Verdict, fmt.Sprintf("~UNK~[%d]", e.Verdict))
		cost := time.Duration(e.IptDelay)
		iptables = fmt.Sprintf("pf=%s, table=%s hook=%s verdict=%s cost=%s", pf, iptName, hook, verdict, cost)
	}

	funcName := NullTerminatedStringToString(e.FuncName[:])
	pktType := e.outputPktType(e.PktType)
	if iptables == "" {
		return fmt.Sprintf("pkt_type=%s func=%s", pktType, funcName)
	}
	return fmt.Sprintf("pkt_type=%s iptables=[%s]", pktType, iptables)
}

func (e *PerfEvent) outputPktType(pktType uint8) string {
	// See: https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/if_packet.h#L26
	const (
		PACKET_USER   = 6
		PACKET_KERNEL = 7
	)
	pktTypes := map[uint8]string{
		syscall.PACKET_HOST:      "HOST",
		syscall.PACKET_BROADCAST: "BROADCAST",
		syscall.PACKET_MULTICAST: "MULTICAST",
		syscall.PACKET_OTHERHOST: "OTHERHOST",
		syscall.PACKET_OUTGOING:  "OUTGOING",
		syscall.PACKET_LOOPBACK:  "LOOPBACK",
		PACKET_USER:              "USER",
		PACKET_KERNEL:            "KERNEL",
	}
	if s, ok := pktTypes[pktType]; ok {
		return s
	}
	return ""
}

func (e *PerfEvent) Output(ipt *IptablesInfo, trace *IptablesTrace) string {
	var s strings.Builder

	// time
	t := e.outputTimestamp()
	s.WriteString(fmt.Sprintf("[%-8s] ", t))

	// skb
	s.WriteString(fmt.Sprintf("[0x%-16x] ", e.Skb))

	// netns
	s.WriteString(fmt.Sprintf("[%-10d] ", e.NetNS))

	// pid
	s.WriteString(fmt.Sprintf("%-8d ", e.Pid))

	// comm
	comm := nullTerminatedStr(e.Comm[:])
	s.WriteString(fmt.Sprintf("%-12s ", comm))

	// cpu
	s.WriteString(fmt.Sprintf("%-4d ", e.CPU))

	// interface
	ifname := nullTerminatedStr(e.Ifname[:])
	if ifname == "" {
		ifname = "nil"
	}
	s.WriteString(fmt.Sprintf("%-14s ", ifname))

	// EthProto
	//ethProto := tools.EthProtoString(int(tools.Htons(uint16(e.EthProto))))
	//ethProto := tools.EthProtoString(int(e.l2Info.L3Proto))
	//if tools.EthProtoString(int(tools.Htons(e.VlanProto))) == "VLAN" {
	//	ethProto = "VLAN"
	//}
	//s.WriteString(fmt.Sprintf("%-12s ", ethProto))

	// src mac
	srcMac := net.HardwareAddr(e.SrcMac[:]).String()
	s.WriteString(fmt.Sprintf("%-18s ", srcMac))

	// dest mac
	destMac := net.HardwareAddr(e.DestMac[:]).String()
	s.WriteString(fmt.Sprintf("%-18s ", destMac))

	// Nexthop
	nexthop := tools.Mac2Ip(net.HardwareAddr(e.DestMac[:]).String())
	s.WriteString(fmt.Sprintf("%-14s ", nexthop))

	// vlanId
	s.WriteString(fmt.Sprintf("%-8d ", e.VlanTci))

	// ip len
	s.WriteString(fmt.Sprintf("%-8d ", e.TotLen))

	// pkt info
	pktInfo := e.outputPktInfo()
	s.WriteString(fmt.Sprintf("%-46s ", pktInfo))

	// trace info
	//traceInfo := e.outputTraceInfo()
	traceInfo := e.outputIptablesInfo(ipt, trace)
	s.WriteString(traceInfo)

	return s.String()
}
