package plugin

import "github.com/twhiston/factd/lib/common"

// The Version plugin provides information about the factd version number
type Version struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *Version) Name() string {
	return GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *Version) Report(facts chan<- ReportedFact) {
	OneShotReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *Version) Facts() (common.FactList, error) {
	facts := common.FactList{"version": common.FactdVersion}
	return facts, nil
}
