// +build mage

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

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/mysteriumnetwork/go-dvpn-web/ci/command"
)

// Generate re-generates the assets_vfsdata.go
func Generate() error {
	defer command.Cleanup()
	mg.SerialDeps(
		command.DownloadLatestAssets,
		command.ExtractAssets,
		command.FixDirectory,
		command.Generate,
	)
	return nil
}

func Commit() error {
	mg.Deps(command.PushChange)
	return nil
}
