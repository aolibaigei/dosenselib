package hostinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cobaltsense/dosenselib/fileutil"
	"github.com/shirou/gopsutil/host"
)

type InfoStat struct {
	Hostname        string `json:"hostname"`
	Uptime          uint64 `json:"uptime"`
	BootTime        uint64 `json:"bootTime"`
	Procs           int    `json:"procs"`           // number of processes
	OS              string `json:"os"`              // ex: freebsd, linux
	Platform        string `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily  string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion string `json:"platformVersion"` // version of the complete OS
	KernelVersion   string `json:"kernelVersion"`   // version of the OS kernel (if available)
	HostID          string `json:"hostid"`          // ex: uuid
}

func KernelVersion() (string, error) {

	filename := "/proc/sys/kernel/osrelease"
	version := ""
	if fileutil.IsExist(filename) == true {
		data, err := ioutil.ReadFile(filename)
		if err == nil {
			version = strings.Replace(string(data), "\n", "", -1)
		}

	}
	return version, nil
}

func HostID() string {

	filename := "/etc/machine-id"
	machine_id := ""
	if fileutil.IsExist(filename) == true {
		data, err := ioutil.ReadFile(filename)
		if err == nil {
			machine_id = strings.Replace(string(data), "\n", "", -1)
		}

	}

	host_id := fmt.Sprintf("%s-%s-%s-%s-%s", machine_id[0:8], machine_id[8:12], machine_id[12:16], machine_id[16:20], machine_id[20:32])

	return host_id
}

func Info() (*InfoStat, error) {

	ret := &InfoStat{
		OS: runtime.GOOS,
	}

	hostname, err := os.Hostname()
	if err == nil {
		ret.Hostname = hostname
	}

	platform, family, version, err := host.PlatformInformation()
	if err == nil {
		ret.Platform = platform
		ret.PlatformFamily = family
		ret.PlatformVersion = version
	}
	kernelVersion, err := KernelVersion()
	if err == nil {
		ret.KernelVersion = kernelVersion
	}
	boot, err := host.BootTime()
	if err == nil {
		ret.BootTime = boot
		ret.Uptime = uint64(time.Now().Unix()) - boot
	}

	ret.Procs = runtime.NumCPU()
	ret.HostID = HostID()

	return ret, nil
}
