package plugin

import (
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/common/logging"
	"reflect"
	"strings"
	"time"
)

// Plugin interface represents a set of facts grouped under a common name
type Plugin interface {
	// Name returns the plugin name, as a string
	// commonly used as the reference on the cli or in config files
	Name() string

	// Return a map of facts, keyed by string
	// This can be called independently of Report
	Facts() (common.FactList, error)

	// The report function will be called in a goroutine and is expected to pass a ReportedFact to the channel
	// If your report is a one off event and should not wait or poll the data will be persisted to the facts data for the duration of the run
	// Usually Report is a wrapper for Fact, creating an object it can write to a channel
	// A standard implementation of this would be
	// func (p *MyPlugin)Report(facts chan<- plugin.ReportedFact){
	// 	  plugin.PollingReport(p,facts)
	// }
	Report(facts chan<- ReportedFact)
}

// ReportedFact represents a fact value
// Because the reported fact is decoupled from its (potential) parent you need to specify the parent key in the fact report
// The interface type returned must match the interface type returned by the plugin Facts() method
type ReportedFact struct {
	// Top level key associated with this fact. Usually Plugin.Name()
	Parent string
	// A list of facts to be reported
	Facts common.FactList
}

// GetPluginName is a helper function to return the class name nicely formatted
// for use in the Name function
func GetPluginName(myvar interface{}) string {
	name := reflect.TypeOf(myvar).String()
	parts := strings.Split(name, ".")
	return strings.ToLower(parts[len(parts)-1])
}

// PollingReport is a helper function for providing polling report data
// TODO - make poll time an option
func PollingReport(p Plugin, facts chan<- ReportedFact) {

	ticker := time.NewTicker(time.Second * time.Duration(5))
	for { //nolint: megacheck
		select {
		case <-ticker.C:
			fact := ReportedFact{Parent: p.Name()}
			data, err := p.Facts()
			if err != nil {
				logging.HandleError(err)
			}
			fact.Facts = data
			facts <- fact
		}
	}
}

// OneShotReport is a wrapper for getting facts and wrapping them in a reported fact
// usually called in a plugin from the Report function
func OneShotReport(p Plugin, facts chan<- ReportedFact) {

	fact := ReportedFact{Parent: p.Name()}
	data, err := p.Facts()
	logging.HandleError(err)
	fact.Facts = data
	facts <- fact

}
