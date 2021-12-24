/*
 * Copyright (C) 2019 The "MysteriumNetwork/go-dvpn-web" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mholt/archiver/v3"
)

const localDvpnWebPath = "DVPN_WEB_LOCAL_PATH"
const assetName = "dist.tar.gz"
const tempDir = "temp"
const assetDir = "assets"

// GenerateLocal hello world
func GenerateLocal() error {
	if os.Getenv(localDvpnWebPath) == "" {
		return errors.New("please set " + localDvpnWebPath + " environment variable")
	}

	defer Cleanup()
	mg.SerialDeps(
		CopyLocalAssets,
		ExtractAssets,
		FixDirectory,
		GoGenerate,
	)
	return nil
}

func CopyLocalAssets() error {
	pathToDist := os.Getenv(localDvpnWebPath) + "/" + assetName
	fmt.Println("Copying assets from: " + pathToDist)
	out, err := os.Create(assetName)
	if err != nil {
		return err
	}
	defer out.Close()

	open, _ := os.Open(pathToDist)
	_, err = io.Copy(out, open)
	defer open.Close()
	return err
}

// Generate re-generates the assets_vfsdata.go
func Generate() error {
	defer Cleanup()
	mg.SerialDeps(
		DownloadLatestAssets,
		ExtractAssets,
		FixDirectory,
		GoGenerate,
	)
	return nil
}

// DownloadLatestAssets fetches the latest assets from github
func DownloadLatestAssets() error {
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
		return errors.New("no assets in latest release")
	}

	found := -1
	for i, v := range rr.Assets {
		if v.Name == assetName {
			found = i
		}
	}

	if found < 0 {
		return fmt.Errorf("no %q found in assets of release", assetName)
	}

	fmt.Println("downloading file", rr.Assets[found].BrowserDownloadURL)
	err = downloadFile(assetName, rr.Assets[found].BrowserDownloadURL)
	if err != nil {
		return err
	}

	fmt.Println("file downloaded")
	return nil
}

// ExtractAssets extracts the asset
func ExtractAssets() error {
	z := archiver.TarGz{
		Tar: &archiver.Tar{
			OverwriteExisting: true,
			MkdirAll:          true,
		},
	}
	fmt.Println("extracting archive", assetName)
	err := z.Unarchive(assetName, tempDir)
	if err != nil {
		return err
	}
	fmt.Println("archive", assetName, "extracted")
	return nil
}

// FixDirectory unnests the assets
func FixDirectory() error {
	fmt.Println("renaming directory assets/build -> assets")
	err := os.Rename(tempDir+"/build", assetDir)
	if err != nil {
		return err
	}
	fmt.Println("directory renamed")
	return nil
}

// Cleanup removes the assets
func Cleanup() error {
	fmt.Println("cleaning up...")
	toClean := []string{
		tempDir, assetDir, assetName,
	}

	for _, v := range toClean {
		err := os.RemoveAll(v)
		if err != nil {
			fmt.Println("could not remove", v)
		}
	}
	fmt.Println("cleanup done")
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

// Generate generates an updated assets_vfsdata.go file
func GoGenerate() error {
	return sh.RunV("go", "generate")
}

// ReleasesResponse represents the github releases response
type ReleasesResponse struct {
	URL             string      `json:"url"`
	AssetsURL       string      `json:"assets_url"`
	UploadURL       string      `json:"upload_url"`
	HTMLURL         string      `json:"html_url"`
	ID              int         `json:"id"`
	NodeID          string      `json:"node_id"`
	TagName         string      `json:"tag_name"`
	TargetCommitish string      `json:"target_commitish"`
	Name            interface{} `json:"name"`
	Draft           bool        `json:"draft"`
	Author          Author      `json:"author"`
	Prerelease      bool        `json:"prerelease"`
	CreatedAt       time.Time   `json:"created_at"`
	PublishedAt     time.Time   `json:"published_at"`
	Assets          []Assets    `json:"assets"`
	TarballURL      string      `json:"tarball_url"`
	ZipballURL      string      `json:"zipball_url"`
	Body            interface{} `json:"body"`
}

// Author represents the release author
type Author struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// Uploader represents the artifact uploader
type Uploader struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// Assets represents the asset
type Assets struct {
	URL                string    `json:"url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	Uploader           Uploader  `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}
