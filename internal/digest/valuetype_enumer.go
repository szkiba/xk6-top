// Code generated by "enumer -text -json -transform lower -trimprefix ValueType -type ValueType"; DO NOT EDIT.

package digest

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _ValueTypeName = "defaulttimedata"

var _ValueTypeIndex = [...]uint8{0, 7, 11, 15}

const _ValueTypeLowerName = "defaulttimedata"

func (i ValueType) String() string {
	if i < 0 || i >= ValueType(len(_ValueTypeIndex)-1) {
		return fmt.Sprintf("ValueType(%d)", i)
	}
	return _ValueTypeName[_ValueTypeIndex[i]:_ValueTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ValueTypeNoOp() {
	var x [1]struct{}
	_ = x[ValueTypeDefault-(0)]
	_ = x[ValueTypeTime-(1)]
	_ = x[ValueTypeData-(2)]
}

var _ValueTypeValues = []ValueType{ValueTypeDefault, ValueTypeTime, ValueTypeData}

var _ValueTypeNameToValueMap = map[string]ValueType{
	_ValueTypeName[0:7]:        ValueTypeDefault,
	_ValueTypeLowerName[0:7]:   ValueTypeDefault,
	_ValueTypeName[7:11]:       ValueTypeTime,
	_ValueTypeLowerName[7:11]:  ValueTypeTime,
	_ValueTypeName[11:15]:      ValueTypeData,
	_ValueTypeLowerName[11:15]: ValueTypeData,
}

var _ValueTypeNames = []string{
	_ValueTypeName[0:7],
	_ValueTypeName[7:11],
	_ValueTypeName[11:15],
}

// ValueTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ValueTypeString(s string) (ValueType, error) {
	if val, ok := _ValueTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ValueTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ValueType values", s)
}

// ValueTypeValues returns all values of the enum
func ValueTypeValues() []ValueType {
	return _ValueTypeValues
}

// ValueTypeStrings returns a slice of all String values of the enum
func ValueTypeStrings() []string {
	strs := make([]string, len(_ValueTypeNames))
	copy(strs, _ValueTypeNames)
	return strs
}

// IsAValueType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ValueType) IsAValueType() bool {
	for _, v := range _ValueTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ValueType
func (i ValueType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ValueType
func (i *ValueType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ValueType should be a string, got %s", data)
	}

	var err error
	*i, err = ValueTypeString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for ValueType
func (i ValueType) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for ValueType
func (i *ValueType) UnmarshalText(text []byte) error {
	var err error
	*i, err = ValueTypeString(string(text))
	return err
}