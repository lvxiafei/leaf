package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"strings"

	"github.com/blang/semver/v4"
)

func isKernelVersionGte_5_16() (bool, error) {
	//release, err := host.KernelVersion()
	//if err != nil {
	//	return false, fmt.Errorf("failed to get kernel version: %v", err)
	//}

	//release := "4.18.0-425.3.1.el8.x86_64"
	release := "5.14.0-284.25.1.el9_2.x86_64"
	log.Printf("kernel version: %v", release)

	release = strings.Replace(release, ".x86_64", "", -1)

	//if release == "4.18.0-425.3.1.el8.x86_64" {
	//	release = "4.18.0-425.3.1.el8"
	//}

	version, err := semver.Make(release)
	if err != nil {
		return false, fmt.Errorf("failed to parse kernel version: %v", err)
	}

	version_5_16 := semver.MustParse("5.16.0")

	return version.GTE(version_5_16), nil
}

func isKernelVersionGte_5_16_0() (bool, error) {
	release, err := host.KernelVersion()
	if err != nil {
		return false, fmt.Errorf("failed to get kernel version: %v", err)
	}
	//release = "5.14.0-284.25.1.el9_2.x86_64"
	release = "5.4.0-122-generic"
	if true {
		log.Printf("kernel version: %v", release)
	}

	release = release[:6]
	release = strings.Trim(release, "-")
	log.Printf("kernel version part: %v", release)

	version, err := semver.Make(release)
	if err != nil {
		return false, fmt.Errorf("failed to parse kernel version: %v", err)
	}

	version_5_16 := semver.MustParse("5.16.0")

	return version.GTE(version_5_16), nil
}
func main() {
	isHighVersion, err := isKernelVersionGte_5_16_0()
	if err != nil {
		log.Printf("Failed to check kernel version: %v", err)
		return
	}
	if isHighVersion {
		log.Printf("isHighVersion...")
	} else {
		log.Printf("not isHighVersion...")
	}
}
