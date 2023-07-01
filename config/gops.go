package config

import (
	"fmt"
	"os"

	"github.com/google/gops/agent"
)

func RunGops() {
	if Cfg.Gops == "" {
		return
	}

	listenAddr := Cfg.Gops
	options := agent.Options{
		Addr: listenAddr,
	}
	if err := agent.Listen(options); err != nil {
		fmt.Fprintln(os.Stderr, "failed to start gops with addr:", listenAddr)
	} else {
		fmt.Fprintln(os.Stdout, "gops is listening on", listenAddr)
	}
}
