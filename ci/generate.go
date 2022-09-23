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
	"github.com/magefile/mage/mg"
)

const distAssetName = "dist.tar.gz"
const compatibilityAssetName = "compatibility.json"
const versionAssetName = "version.json"
const tempDir = "temp"
const assetDir = "assets"

// Generate re-generates the assets_vfsdata.go
func Generate() error {
	defer Cleanup()
	mg.SerialDeps(
		DownloadAssets,
		ExtractAssets,
		FixDirectory,
		GoGenerate,
		GenerateNodeUIVersionManifest,
	)
	return nil
}
