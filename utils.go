package main

import "strings"

func getNumberTypes() []string {
	return []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "rune"}
}

func getFloatTypes() []string {
	return []string{"float32", "float64"}
}

func getNumberTypeAliases() []string {
	return []string{"uintptr", "complex64"}
}

func getBooleanTypes() []string {
	return []string{"bool"}
}

func getStringTypes() []string {
	return []string{"string"}
}

func getTsTypeName(typeName string) string {
	lowered := strings.ToLower(typeName)
	if SliceContains(getNumberTypes(), lowered) {
		return "number"
	} else if SliceContains(getFloatTypes(), lowered) {
		return "number"
	} else if SliceContains(getNumberTypeAliases(), lowered) {
		return "number"
	} else if SliceContains(getBooleanTypes(), lowered) {
		return "boolean"
	} else if SliceContains(getStringTypes(), lowered) {
		return "string"
	}
	return typeName
}

func SliceContains(slice []string, item string) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	return false
}
