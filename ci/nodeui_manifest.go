package ci

import (
	"encoding/json"
	"errors"
	"fmt"
	dvpnweb "github.com/mysteriumnetwork/go-dvpn-web/v2"
	"io/ioutil"
	"os"
)

func GenerateNodeUIVersionManifest() error {
	tagVersion := os.Getenv("GIT_TAG_VERSION")
	if tagVersion == "" {
		return errors.New("please specify the GIT_TAG_VERSION environment variable")
	}

	jsonBytes, err := json.Marshal(dvpnweb.UIVersion{Version: tagVersion})
	if err != nil {
		return fmt.Errorf("could not marshal version: %w", err)
	}

	err = ioutil.WriteFile(versionAssetName, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("could write to "+versionAssetName+": %w", err)
	}

	return nil
}
