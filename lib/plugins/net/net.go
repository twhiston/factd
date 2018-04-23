package net

import (
	"github.com/fatih/structs"
	n "github.com/shirou/gopsutil/net"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugins"
	"net"
	"strings"
)

// The Net plugin provides information about networks
type Net struct{}

// Name returns the plugins printable name, also used as the map key in the master fact list
func (p *Net) Name() string {
	return plugins.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *Net) Report(facts chan<- plugins.ReportedFact) {
	plugins.PollingReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *Net) Facts() (common.FactList, error) {

	data := make(common.FactList)

	netIfaces, err := n.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, v := range netIfaces {
		mapped := structs.Map(v)
		mapped["IP"] = splitAddrByType(v)
		data[v.Name] = mapped
	}

	return data, nil
}

// splitAddrByType returns a map keyed by ip type v4/v6 containing the appropriate ip's for the interface
func splitAddrByType(v n.InterfaceStat) map[string][]string {
	addrsRep := make(map[string][]string)
	addrsRep["v4"] = make([]string, 0)
	addrsRep["v6"] = make([]string, 0)
	for _, ad := range v.Addrs {
		parts := strings.Split(ad.Addr, "/")
		parsedIP := net.ParseIP(parts[0])
		aType := "v4"
		if parsedIP.To4() == nil {
			aType = "v6"
		}
		addrsRep[aType] = append(addrsRep[aType], ad.Addr)
	}
	return addrsRep
}
