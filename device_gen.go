// Code generated by "enumer -type=Device -output=device_gen.go -transform=snake"; DO NOT EDIT.

package homlet

import (
	"fmt"
)

const _DeviceName = "unknownjeenode_thlmjeenode_thltinytx_ttinytx_thtinytx_tl"

var _DeviceIndex = [...]uint8{0, 7, 19, 30, 38, 47, 56}

func (i Device) String() string {
	if i < 0 || i >= Device(len(_DeviceIndex)-1) {
		return fmt.Sprintf("Device(%d)", i)
	}
	return _DeviceName[_DeviceIndex[i]:_DeviceIndex[i+1]]
}

var _DeviceValues = []Device{0, 1, 2, 3, 4, 5}

var _DeviceNameToValueMap = map[string]Device{
	_DeviceName[0:7]:   0,
	_DeviceName[7:19]:  1,
	_DeviceName[19:30]: 2,
	_DeviceName[30:38]: 3,
	_DeviceName[38:47]: 4,
	_DeviceName[47:56]: 5,
}

var _DeviceNames = []string{
	_DeviceName[0:7],
	_DeviceName[7:19],
	_DeviceName[19:30],
	_DeviceName[30:38],
	_DeviceName[38:47],
	_DeviceName[47:56],
}

// DeviceString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DeviceString(s string) (Device, error) {
	if val, ok := _DeviceNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Device values", s)
}

// DeviceValues returns all values of the enum
func DeviceValues() []Device {
	return _DeviceValues
}

// DeviceStrings returns a slice of all String values of the enum
func DeviceStrings() []string {
	strs := make([]string, len(_DeviceNames))
	copy(strs, _DeviceNames)
	return strs
}

// IsADevice returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Device) IsADevice() bool {
	for _, v := range _DeviceValues {
		if i == v {
			return true
		}
	}
	return false
}
