package mem

import (
	"fmt"

	"github.com/fatih/structs"
	m "github.com/shirou/gopsutil/mem"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/common/logging"
	"github.com/twhiston/factd/lib/plugin"
)

// The Mem plugin provides information about memory limits and usage
type Mem struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *Mem) Name() string {
	return plugin.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *Mem) Report(facts chan<- plugin.ReportedFact) {
	plugin.PollingReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *Mem) Facts() (common.FactList, error) {

	data := make(common.FactList)

	memory, err := getMemoryData()
	logging.HandleError(err)
	data["memory"] = memory

	swap, err := getSwapData()
	logging.HandleError(err)
	data["swap"] = swap

	return data, nil
}

func getMemoryData() (map[string]interface{}, error) {
	memory := make(map[string]interface{})

	// Get the virtual memory from gopsutil
	hostVMem, err := m.VirtualMemory()
	if err != nil {
		return nil, err
	}
	memory["raw"] = structs.Map(hostVMem)

	// Add a memoryfree facts
	free := make(map[string]interface{})
	err = addMemoryUnits(&free, hostVMem.Free)
	if err != nil {
		return memory, err
	}
	memory["free"] = free

	// Add a memorytotal facts
	total := make(map[string]interface{})
	err = addMemoryUnits(&total, hostVMem.Total)
	if err != nil {
		return memory, err
	}
	memory["total"] = total

	return memory, err
}
func getSwapData() (map[string]interface{}, error) {

	swap := make(map[string]interface{})
	// Get the swap information from gopsutil
	hostSwapMem, err := m.SwapMemory()
	if err != nil {
		return nil, err
	}
	swap["raw"] = structs.Map(hostSwapMem)

	// Add a memoryfree facts
	free := make(map[string]interface{})
	err = addMemoryUnits(&free, hostSwapMem.Free)
	if err != nil {
		return swap, err
	}
	swap["free"] = free

	// Add a memorytotal facts
	total := make(map[string]interface{})
	err = addMemoryUnits(&total, hostSwapMem.Total)
	if err != nil {
		return swap, err
	}
	swap["total"] = total

	return swap, err
}

// addMemoryUnits will convert a memory fact into "fact_mb" and "fact_bytes"
func addMemoryUnits(data *map[string]interface{}, memory uint64) error {
	units := map[string]string{
		"GB": "gb",
		"MB": "mb",
		"B":  "bytes",
	}

	for unit, unitLabel := range units {
		convertedMemory, _, err := common.ConvertBytesTo(memory, unit)
		if err != nil {
			return err
		}

		(*data)[unitLabel] = fmt.Sprintf("%.2f", convertedMemory)
	}

	return nil
}
