package types

import "shake/bimap"

type Type int

const (
	TypeEmpty Type = iota
	TypeInt32
	TypeInt64
	TypeUnknown
)

var typeNamesMap = map[Type]string{
	TypeEmpty:   "empty",
	TypeInt32:   "int32",
	TypeInt64:   "int64",
	TypeUnknown: "unknown",
}

// Create the bidirectional map from the map
var typeNames, err = bimap.NewBiMapFromMap(typeNamesMap)

func (t Type) String() string {
	if err != nil {
		panic("Failed to initialize typeNames: " + err.Error())
	}

	// Check if the type exists in the bidirectional map
	value, ok := typeNames.GetByKey(t)
	if !ok {
		return "Unknown"
	}
	return value
}

func GetType(typeName string) Type {
	if err != nil {
		panic("Failed to initialize typeNames: " + err.Error())
	}

	// Check if the typeName exists in the bidirectional map
	value, ok := typeNames.GetByValue(typeName)
	if !ok {
		return TypeUnknown
	}
	return value
}
