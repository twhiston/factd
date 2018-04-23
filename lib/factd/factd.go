package factd

import (
	"bytes"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugins"

	"github.com/twhiston/factd/lib/common/logging"
)

// Factd struct holds Factd-related attributes
type Factd struct {
	facts         map[string]common.FactList
	config        Config
	reportChannel chan plugins.ReportedFact
}

// New returns new instance of Factd
func New(userConf Config) *Factd {
	f := &Factd{
		facts:         make(map[string]common.FactList),
		config:        userConf,
		reportChannel: make(chan plugins.ReportedFact),
		//reportStopChannel: make(chan struct{}),
	}
	return f
}

// MergeSet adds a fact set, merging it with existing fact data
func (f *Factd) MergeSet(k string, v common.FactList) {
	f.ensureFact(k)
	src := f.facts[k]
	err := mergo.Merge(&src, v)
	logging.HandleError(err)
	f.facts[k] = src
}

// ReplaceSet replaces all facts for key k
func (f *Factd) ReplaceSet(k string, v common.FactList) {
	f.ensureFact(k)
	f.facts[k] = v
}

// AddFact sets fact in group at key id
func (f *Factd) AddFact(group string, id string, fact interface{}) {
	f.ensureFact(group)
	f.facts[group][id] = fact
}

// ensureFact makes a factList for a name if it doesnt exist already
func (f *Factd) ensureFact(name string) {
	if f.facts[name] == nil {
		f.facts[name] = make(common.FactList)
	}
}

// Delete deletes given fact
func (f *Factd) Delete(k string) {
	delete(f.facts, k)
}

// Get returns value of given fact, if it exists
func (f *Factd) Get(group string, k string) (interface{}, bool) {
	value, ok := f.facts[group][k]
	return value, ok
}

// Format returns the formatted result of the facts as a pointer to a bytes.Buffer
func (f *Factd) Format() (*bytes.Buffer, error) {
	formatted, err := f.config.Formatter.Format(f.facts)
	if err != nil {
		return &bytes.Buffer{}, err
	}
	return formatted, nil
}

// Print facts by calling selected formatter
// and printing returned buffer to f.Config.Output as a string
func (f *Factd) Print() error {
	formatted, err := f.Format()
	if err != nil {
		return err
	}
	fmt.Fprint(*f.config.Output, formatted.String())
	return nil
}

// GetConfig returns the current config for factd
func (f *Factd) GetConfig() Config {
	return f.config
}

// SetConfig sets the factd config to the provided struct
func (f *Factd) SetConfig(c Config) {
	f.config = c
}

// Collect facts from each active plugin and add them to the plugin result array
func (f *Factd) Collect() {
	for _, v := range f.config.Plugins {
		facts, err := v.Facts()
		logging.HandleError(err)
		f.MergeSet(v.Name(), facts)
	}
}

// RunReporters starts all of the activeReporters in goroutines
// It additionally starts the reporter processor which will listen
// for channel input and merge it with existing reports
func (f *Factd) RunReporters() {
	//TODO - how to stop this being run repeatedly
	for _, v := range f.config.Plugins {
		go v.Report(f.reportChannel)
	}
	go f.ReporterProcessor()

}

// StopReporters will close the report channel
func (f *Factd) StopReporters() {
	//TODO - needs to be a bit smarter
	//close(f.reportStopChannel)
	close(f.reportChannel)
}

// ReporterProcessor is responsible for acting on information passed to the report channel
// it will trigger the listeners with the provided fact BEFORE it merges the new fact into the parent set
func (f *Factd) ReporterProcessor() {
	for { //nolint: megacheck
		select {
		case fact := <-f.reportChannel:
			f.MergeSet(fact.Parent, fact.Facts)
		}
	}
}
