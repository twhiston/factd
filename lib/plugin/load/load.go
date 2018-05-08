package load

import (
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugin"

	"github.com/fatih/structs"
	l "github.com/shirou/gopsutil/load"
)

// The Load plugin provides information about current load on the server
type Load struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *Load) Name() string {
	return plugin.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *Load) Report(facts chan<- plugin.ReportedFact) {
	plugin.PollingReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *Load) Facts() (common.FactList, error) {

	data := make(common.FactList)

	avg, err := l.Avg()
	if err != nil {
		return nil, err
	}
	data["Avg"] = structs.Map(avg)
	misc, err := l.Misc()
	if err != nil {
		return data, err
	}
	data["Misc"] = structs.Map(misc)
	return data, nil
}
