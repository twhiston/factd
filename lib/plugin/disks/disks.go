package disks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"

	"strconv"

	d "github.com/shirou/gopsutil/disk"
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugin"
)

var reDevBlacklist = regexp.MustCompile("^(dm-[0-9]+|loop[0-9]+)$")

// The Disks plugin provides information about disks and partitions of the machine
type Disks struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *Disks) Name() string {
	return plugin.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *Disks) Report(facts chan<- plugin.ReportedFact) {
	plugin.PollingReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *Disks) Facts() (common.FactList, error) {

	data := make(common.FactList)

	partitions, err := d.Partitions(false)
	if err != nil {
		return nil, err
	}
	data["partitions"] = partitions

	serials := make(map[string]string)
	for _, v := range partitions {
		s := d.GetDiskSerialNumber(v.Device)
		if s != "" {
			serials[v.Device] = s
		}
	}
	data["serials"] = serials

	//TODO - make this plugin config
	blockDevs, err := getBlockDevices(false)
	if err != nil {
		return nil, err
	}
	data["devices"] = blockDevs

	return data, nil
}

// getBlockDevices returns list of block devices
func getBlockDevices(all bool) (map[string]map[string]string, error) {
	blockDevs := make(map[string]map[string]string)
	targetDir := fmt.Sprintf("%v/block", common.GetHostSys())
	contents, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return blockDevs, err
	}
	for _, v := range contents {
		if !all {
			if reDevBlacklist.MatchString(v.Name()) {
				continue
			}
		}
		blockDevs[v.Name()] = make(map[string]string)
	}

	for k := range blockDevs {
		size, err := getBlockDeviceSize(k)
		if err == nil {
			blockDevs[k]["size"] = strconv.Itoa(int(size))
		}

		model, err := getBlockDeviceModel(k)
		if err == nil {
			blockDevs[k]["model"] = model
		}

		vendor, err := getBlockDeviceVendor(k)
		if err == nil {
			blockDevs[k]["vendor"] = vendor
		}
	}

	return blockDevs, nil
}

// getBlockDeviceModel returns model of block device as reported by Linux
// kernel.
func getBlockDeviceModel(blockDevice string) (string, error) {
	modelFilename := fmt.Sprintf("%s/block/%s/device/model",
		common.GetHostSys(), blockDevice)
	model, err := ioutil.ReadFile(modelFilename)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bytes.TrimSuffix(model, []byte("\n"))), nil
}

// getBlockDeviceSize returns size of block device as reported by Linux kernel
// multiplied by 512.
func getBlockDeviceSize(blockDevice string) (int64, error) {
	sizeFilename := fmt.Sprintf("%s/block/%s/size", common.GetHostSys(),
		blockDevice)
	size, err := ioutil.ReadFile(sizeFilename)
	if err != nil {
		return 0, err
	}
	sizeInt, err := strconv.ParseInt(fmt.Sprintf("%s",
		bytes.TrimSuffix(size, []byte("\n"))), 10, 64)
	if err != nil {
		return 0, err
	}
	return sizeInt * 512, nil
}

// getBlockDeviceVendor returns vendor of block device as reported by Linux
// kernel.
func getBlockDeviceVendor(blockDevice string) (string, error) {
	vendorFilename := fmt.Sprintf("%s/block/%s/device/vendor",
		common.GetHostSys(), blockDevice)
	vendor, err := ioutil.ReadFile(vendorFilename)
	if err != nil {
		return "", err
	}
	vendor = bytes.TrimSuffix(vendor, []byte("\n"))
	vendor = bytes.TrimRight(vendor, " ")
	return fmt.Sprintf("%s", vendor), nil
}
