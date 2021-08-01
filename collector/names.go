package collector

import (
	"strconv"
	"strings"
)

func cleanupName(name string) string {
	switch name {
	case "type", "struct", "map", "chan", "select", "switch", "case", "default", "range", "func",
		"package", "import", "var", "const", "return":
		return cleanupName(name + "_")
	default:
		return name
	}
}

func publicName(name string) string {
	name = strings.TrimPrefix(name, "_")
	// Handle cases like DXGI_MATRIX_3X2_F fields.
	if _, err := strconv.Atoi(name); err == nil {
		name = "Field_" + name
	}
	return strings.Title(cleanupName(name))
}
