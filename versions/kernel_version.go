package versions

import (
	"fmt"
	"log"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/shirou/gopsutil/v3/host"
)

func IsKernelVersionGte_5_16_0() (bool, error) {
	release, err := host.KernelVersion()
	if err != nil {
		return false, fmt.Errorf("failed to get kernel version: %v", err)
	}
	
	release = strings.Trim(release[:6], "-")
	log.Printf("kernel version part: %v", release)

	version, err := semver.Make(release)
	if err != nil {
		return false, fmt.Errorf("failed to parse kernel version: %v", err)
	}

	version_5_16 := semver.MustParse("5.16.0")

	return version.GTE(version_5_16), nil
}
