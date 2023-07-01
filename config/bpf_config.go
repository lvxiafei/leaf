package config

type BpfConfig struct {
	NetNS  uint32
	Pid    uint32
	IP     uint32
	Port   uint16
	IcmpID uint16
	Proto  uint8
	Pad    [3]uint8
}

func GetBpfConfig() BpfConfig {
	return BpfConfig{
		NetNS:  Cfg.NetNS,
		Pid:    Cfg.Pid,
		IP:     Cfg.ip,
		Port:   (Cfg.Port >> 8) & (Cfg.Port << 8),
		IcmpID: (Cfg.IcmpID >> 8) & (Cfg.IcmpID << 8),
		Proto:  Cfg.proto,
	}
}
