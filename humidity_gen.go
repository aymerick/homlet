// Code generated by "enumer -type=HumidityStatus -output=humidity_gen.go -transform=snake -trimprefix=Humidity"; DO NOT EDIT.

package homlet

import (
	"fmt"
)

const _HumidityStatusName = "normalcomfortabledrywet"

var _HumidityStatusIndex = [...]uint8{0, 6, 17, 20, 23}

func (i HumidityStatus) String() string {
	if i < 0 || i >= HumidityStatus(len(_HumidityStatusIndex)-1) {
		return fmt.Sprintf("HumidityStatus(%d)", i)
	}
	return _HumidityStatusName[_HumidityStatusIndex[i]:_HumidityStatusIndex[i+1]]
}

var _HumidityStatusValues = []HumidityStatus{0, 1, 2, 3}

var _HumidityStatusNameToValueMap = map[string]HumidityStatus{
	_HumidityStatusName[0:6]:   0,
	_HumidityStatusName[6:17]:  1,
	_HumidityStatusName[17:20]: 2,
	_HumidityStatusName[20:23]: 3,
}

var _HumidityStatusNames = []string{
	_HumidityStatusName[0:6],
	_HumidityStatusName[6:17],
	_HumidityStatusName[17:20],
	_HumidityStatusName[20:23],
}

// HumidityStatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func HumidityStatusString(s string) (HumidityStatus, error) {
	if val, ok := _HumidityStatusNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to HumidityStatus values", s)
}

// HumidityStatusValues returns all values of the enum
func HumidityStatusValues() []HumidityStatus {
	return _HumidityStatusValues
}

// HumidityStatusStrings returns a slice of all String values of the enum
func HumidityStatusStrings() []string {
	strs := make([]string, len(_HumidityStatusNames))
	copy(strs, _HumidityStatusNames)
	return strs
}

// IsAHumidityStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i HumidityStatus) IsAHumidityStatus() bool {
	for _, v := range _HumidityStatusValues {
		if i == v {
			return true
		}
	}
	return false
}
