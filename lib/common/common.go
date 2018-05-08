package common

import (
	"fmt"
	"math"
	"os"
)

// FactdVersion will be set via init.
// In the default command implementation it is passed from the build process
var FactdVersion = "unknown"

// FactList is a map of data returned from a plugin to be printed by a formatter
type FactList map[string]interface{}

var (
	// ByteUnits is a k=>v map of units for conversion
	ByteUnits = map[int]string{
		0: "B",
		1: "kB",
		2: "MB",
		3: "GB",
		4: "TB",
	}
)

// ConvertBytes converts bytes to the highest possible unit
func ConvertBytes(in uint64) (float64, string, error) {
	out := float64(in)
	idx := 0
	maxIdx := len(ByteUnits)
	for idx < maxIdx && out > 1 {
		tmp := out / 1024
		if tmp < 1 {
			break
		}
		out = tmp
		idx++
	}
	return out, ByteUnits[idx], nil
}

// ConvertBytesTo converts bytes to the specified unit
func ConvertBytesTo(in uint64, maxUnit string) (float64, string, error) {
	if maxUnit == "" {
		return 0, "", fmt.Errorf("invald maximum unit")
	}
	out := float64(in)
	idx := 0
	maxIdx := len(ByteUnits)
	for idx < maxIdx && maxUnit != ByteUnits[idx] {
		out = out / 1024
		idx++
	}
	return out, ByteUnits[idx], nil
}

// ConvertNetmask converts CIDR (netmask) to Netmask
func ConvertNetmask(in uint8) (string, error) {
	if in > 32 {
		return "", fmt.Errorf("invalid netmask")
	}
	octets := map[uint8]uint8{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
	}
	var idx uint8 = 1
	for in > 0 && idx < 5 {
		if (in / 8) > 0 {
			in = in - 8
			octets[idx] = 255
		} else {
			mod := in % 8
			octets[idx] = 255 - uint8(math.Pow(2, float64(8-mod))) + 1
			in = 0
		}
		idx++
	}
	return fmt.Sprintf("%d.%d.%d.%d", octets[1], octets[2], octets[3],
		octets[4]), nil
}

//TODO - this is plugin config for disks
func GetHostSys() string {
	hostSys := os.Getenv("HOST_SYS")
	if hostSys == "" {
		hostSys = "/sys"
	}
	return hostSys
}
