package config

import (
	"encoding/binary"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"leaf/versions"
	"log"
	"net"
	"os"
)

type Debug struct {
	Default           bool `yaml:"default"`
	IndexIptablesInfo bool `yaml:"index_iptables_info"`
}

// Config is the configurations for the bpf program.
type Config struct {
	Hook         map[string]bool `yaml:"hook"`
	Draw         bool            `yaml:"draw"`
	Debug        Debug           `yaml:"debug"`
	File         string          `yaml:"file"`
	Version      bool            `yaml:"version"`
	CatchCount   uint64          `yaml:"catch_count"`
	IP           string          `yaml:"IP"`
	ip           uint32          `yaml:"ip"`
	Proto        string          `yaml:"Proto"`
	Mode         string          `yaml:"mode"`
	Cni          string          `yaml:"cni"`
	Tech         string          `yaml:"tech"`
	Scene        string          `yaml:"scene"`
	proto        uint8           `yaml:"proto"`
	IcmpID       uint16          `yaml:"icmp_id"`
	Port         uint16          `yaml:"port"`
	Pid          uint32          `yaml:"pid"`
	NetNS        uint32          `yaml:"net_ns"`
	Time         bool            `yaml:"time"`
	Timestamp    bool            `yaml:"timestamp"`
	PerCPUBuffer int             `yaml:"per_cpu_buffer"`
	Gops         string          `yaml:"gops"`
}

var Cfg Config

func SetFlags(RootCmd *cobra.Command) {
	fs := RootCmd.PersistentFlags()

	fs.StringVarP(&Cfg.File, "file", "f", "", "config file")
	fs.BoolVar(&Cfg.Debug.Default, "debug", false, "debug info")

	fs.BoolVarP(&Cfg.Draw, "draw", "d", false, "draw graph, default generate")
	fs.StringVarP(&Cfg.IP, "ipaddr", "H", "", "ip address")
	fs.StringVar(&Cfg.Proto, "proto", "", "tcp|udp|icmp|any")
	fs.StringVar(&Cfg.Mode, "mode", "", "generate|define|diff")
	fs.StringVar(&Cfg.Cni, "cni", "", "cilium|kube-ovn|flannel|any")
	fs.StringVar(&Cfg.Tech, "tech", "", "macvlan|vlan|any")
	fs.StringVar(&Cfg.Scene, "scene", "", "node2eip|pod2eip|pod2bms|any")
	fs.Uint16Var(&Cfg.IcmpID, "icmpid", 0, "trace icmp id")
	fs.Uint64VarP(&Cfg.CatchCount, "catch-count", "c", 0, "catch and print count")
	fs.Uint16VarP(&Cfg.Port, "port", "P", 0, "udp or tcp port")
	fs.Uint32VarP(&Cfg.Pid, "pid", "p", 0, "trace this PID only")
	fs.Uint32VarP(&Cfg.NetNS, "netns", "N", 0, "trace this netns inode only")
	fs.BoolVarP(&Cfg.Time, "time", "T", true, "show HH:MM:SS timestamp")
	fs.BoolVarP(&Cfg.Timestamp, "timestamp", "t", false, "show timestamp in seconds at us resolution")
	fs.IntVarP(&Cfg.PerCPUBuffer, "per-cpu-buffer", "B", 4096, "per CPU buffer to receive perf event")
	fs.StringVar(&Cfg.Gops, "gops", "", "gops address")
	fs.BoolVarP(&Cfg.Version, "version", "V", false, "show version")
}

// InitRawConfig 读参数文件
func InitRawConfig(configPath *string) Config {
	c := Config{}

	b, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v, %v", *configPath, err)
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("Failed to unmarshal yaml %v to Object: %v", string(b), err)
	}

	return c
}

func (c *Config) UseSpecialConfig() {

	b, err := os.ReadFile(c.File)
	if err != nil {
		log.Fatalf("Failed to read config file: %v, %v", c.File, err)
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("Failed to unmarshal yaml %v to Object: %v", string(b), err)
	}
}

func (c *Config) Parse() error {
	if Cfg.Version {
		log.Printf(versions.String())
		return fmt.Errorf("  ")
	}

	if len(c.File) > 0 {
		c.UseSpecialConfig()
	}
	log.Printf("cfg : %+v", c)

	ip := c.IP
	if ip != "" {
		ip := net.ParseIP(ip)
		ip = ip.To4()
		if ip == nil {
			return fmt.Errorf("invalid IPv4 addr(%s)", ip)
		}

		//c.ip = binary.BigEndian.Uint32(ip)
		c.ip = binary.LittleEndian.Uint32(ip)
	}

	proto := c.Proto
	if proto != "" {
		switch proto {
		case "tcp":
			c.proto = 6
		case "udp":
			c.proto = 17
		case "icmp":
			c.proto = 1
		case "any":
		default:
			return fmt.Errorf("invalid proto(%s)", proto)
		}
	}
	return nil
}
