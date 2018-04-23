package user

import (
	"github.com/fatih/structs"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugins"
	"os/user"
)

// The User plugin provides information about the user executing factd
type User struct{}

// Name returns the plugins printable name, also used as the map key in the master fact list
func (p *User) Name() string {
	return plugins.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *User) Report(facts chan<- plugins.ReportedFact) {
	plugins.PollingReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *User) Facts() (common.FactList, error) {

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	data := structs.Map(usr)
	data["Groups"], err = usr.GroupIds()

	return data, err
}
