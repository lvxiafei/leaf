package tools

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"leaf/config"
	"log"
)

func KprobeWrapper(symbol string, prog *ebpf.Program, opts *link.KprobeOptions) {

	if len(config.Cfg.File) > 0 && !config.Cfg.Hook[symbol] {
		return
	}

	if _, err := link.Kprobe(symbol, prog, opts); err != nil {
		log.Printf("Failed to kprobe(%s): %v", symbol, err)
		return
	} else {
		log.Printf("Attached kprobe(%s)", symbol)
	}
}
