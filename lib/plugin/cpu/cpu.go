package cpu

import (
	c "github.com/shirou/gopsutil/cpu"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugin"

	"strconv"
)

// The CPU plugin provides information about the cpu's in the machine
type CPU struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *CPU) Name() string {
	return plugin.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *CPU) Report(facts chan<- plugin.ReportedFact) {
	plugin.OneShotReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *CPU) Facts() (common.FactList, error) {

	data := make(common.FactList)

	totalCount, err := c.Counts(true)
	if err != nil {
		return nil, err
	}
	data["ProcessorCount"] = totalCount

	cpuInfo, err := c.Info()
	if err != nil {
		return nil, err
	}

	data["Info"] = cpuInfo

	physIDs := make(map[uint64]interface{})
	for _, v := range cpuInfo {
		pid, err := strconv.ParseUint(v.PhysicalID, 10, 32)
		if err == nil {
			physIDs[pid] = nil
		}
	}
	data["PhysicalProcessorCount"] = len(physIDs)
	return data, nil
}
