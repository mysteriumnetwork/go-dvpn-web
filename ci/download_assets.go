package ci

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// DownloadAssets fetches the latest assets from github
func DownloadAssets() error {
	tagVersion := os.Getenv("GIT_TAG_VERSION")
	if tagVersion == "" {
		return errors.New("please specify the GIT_TAG_VERSION environment variable")
	}

	fmt.Println(fmt.Sprintf("getting dvpn-web release: %s", tagVersion))

	client := &http.Client{
		Timeout: time.Minute,
	}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/mysteriumnetwork/dvpn-web/releases/tags/"+tagVersion, nil)
	if err != nil {
		return err
	}

	if os.Getenv("GIT_TOKEN") != "" {
		fmt.Println("git token present, will make authorized request")
		req.SetBasicAuth("doesntmatter", os.Getenv("GIT_TOKEN"))
	}

	res, err := client.Do(req)
	fmt.Println("response status", res.StatusCode)

	fmt.Println("reading response body")
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println("unmarshaling response body")
	rr := ReleasesResponse{}
	err = json.Unmarshal(bytes, &rr)
	if err != nil {
		return err
	}

	fmt.Println("looking for assets")
	if len(rr.Assets) == 0 {
		return errors.New(fmt.Sprintf("no assets in release: %s", tagVersion))
	}

	if err := findAndDownloadAsset(rr, distAssetName, true); err != nil {
		return err
	}
	if err := findAndDownloadAsset(rr, compatibilityAssetName, false); err != nil {
		return err
	}

	return nil
}

func findAndDownloadAsset(rr ReleasesResponse, assetName string, required bool) error {
	found := -1
	for i, v := range rr.Assets {
		if v.Name == assetName {
			found = i
		}
	}

	if found < 0 && !required {
		fmt.Println(fmt.Sprintf("asset: %s - not found, but is not required - skipping", assetName))
		return nil
	}

	if found < 0 {
		return fmt.Errorf("no %q found in assets of release", assetName)
	}

	fmt.Println("downloading file", rr.Assets[found].BrowserDownloadURL)
	if err := downloadFile(assetName, rr.Assets[found].BrowserDownloadURL); err != nil {
		return err
	}

	fmt.Println("file downloaded")
	return nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
