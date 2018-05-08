package docker

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/fatih/structs"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugin"
)

// The Docker plugin provides information about docker images and containers on the host
type Docker struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *Docker) Name() string {
	return plugin.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *Docker) Report(facts chan<- plugin.ReportedFact) {
	plugin.PollingReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *Docker) Facts() (common.FactList, error) {

	data := make(common.FactList)

	cli, err := docker.NewEnvClient()
	if err != nil {
		return data, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return data, err
	}

	cs := make(common.FactList)
	for _, container := range containers {
		cs[container.ID] = structs.Map(container)
	}
	data["Containers"] = cs

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return data, err
	}
	is := make(common.FactList)
	for _, image := range images {
		is[image.ID] = structs.Map(image)
	}
	data["Images"] = is

	return data, nil
}
