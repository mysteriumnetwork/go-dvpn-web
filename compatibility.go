package dvpnweb

import (
	_ "embed"
	"encoding/json"
)

//go:embed compatibility.json
var compatibilityJson []byte

func MinimalCompatibleNodeVersion() (string, error) {
	var c Compatibility
	return c.MinimalNodeVersion, json.Unmarshal(compatibilityJson, &c)
}

type Compatibility struct {
	MinimalNodeVersion string `json:"minimalNodeVersion"`
}
