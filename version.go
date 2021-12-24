package dvpnweb

import (
	_ "embed"
	"encoding/json"
)

//go:embed version.json
var uiVersionJson []byte

func Version() (string, error) {
	var v UIVersion
	return v.Version, json.Unmarshal(uiVersionJson, &v)
}

type UIVersion struct {
	Version string `json:"version"`
}
