package host

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/structs"
	h "github.com/shirou/gopsutil/host"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugin"
	"runtime"
)

// The Host plugin provides information about the host and OS
type Host struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *Host) Name() string {
	return plugin.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *Host) Report(facts chan<- plugin.ReportedFact) {
	plugin.PollingReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *Host) Facts() (common.FactList, error) {

	hostInfo, err := h.Info()
	if err != nil {
		return nil, err
	}

	data := structs.Map(hostInfo)

	splitted := strings.SplitN(hostInfo.Hostname, ".", 2)
	var hostname *string
	if len(splitted) > 1 {
		hostname = &splitted[0]
		data["Domain"] = splitted[1]
	} else {
		hostname = &hostInfo.Hostname
	}
	data["Hostname"] = *hostname

	var isVirtual bool
	if hostInfo.VirtualizationRole == "host" {
		isVirtual = false
	} else {
		isVirtual = true
	}
	data["IsVirtual"] = isVirtual

	data["UptimeMinutes"] = hostInfo.Uptime / 60
	data["UptimeHours"] = hostInfo.Uptime / 60 / 60
	data["UptimeDays"] = hostInfo.Uptime / 60 / 60 / 24

	envPath := os.Getenv("PATH")
	if envPath != "" {
		data["Path"] = envPath
	}

	var uname syscall.Utsname
	err = syscall.Uname(&uname)
	if err == nil {
		kernelRelease := int8ToString(uname.Release)
		kernelVersion := strings.Split(kernelRelease, "-")[0]
		kvSplitted := strings.Split(kernelVersion, ".")
		data["KernelRelease"] = kernelRelease
		data["KernelVersion"] = kernelVersion
		data["KernelMajVersion"] = strings.Join(kvSplitted[0:2], ".")

		hardwareModel := int8ToString(uname.Machine)
		data["HardwareModel"] = hardwareModel
		data["Architecture"] = runtime.GOARCH
	}

	z, _ := time.Now().Zone()
	data["Timezone"] = z

	hostid, err := getUniqueID()
	if err == nil {
		data["UniqueID"] = hostid
	}
	return data, nil
}

// getUniqueID returns executes % hostid; and returns its STDOUT as a string.
func getUniqueID() (string, error) {
	// #nosec: Subprocess launching should be audited
	cmd := exec.Command("/usr/bin/hostid")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

// int8ToString converts [65]int8 in syscall.Utsname to string
func int8ToString(bs [65]int8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	return strings.TrimRight(string(b), "\x00")
}
