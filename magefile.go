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
	"errors"
	"fmt"
	"os"

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

// CI generates, commits and pushes the changes to the assets_vfsdata.go
func CI() error {
	gitToken := os.Getenv("GIT_TOKEN")
	if gitToken == "" {
		return errors.New("please specify the GIT_TOKEN environment variable")
	}

	tagVersion := os.Getenv("GIT_TAG_VERSION")
	if tagVersion == "" {
		return errors.New("please specify the TAG_VERSION environment variable")
	}

	git := command.NewCommiter(gitToken)
	err := git.Checkout("master")
	if err != nil {
		return err
	}

	defer command.Cleanup()
	mg.SerialDeps(
		command.DownloadLatestAssets,
		command.ExtractAssets,
		command.FixDirectory,
		command.Generate,
	)

	hash, err := git.Commit(fmt.Sprintf("updating assets_vfsdata.go for %v", tagVersion), "assets_vfsdata.go")
	if err != nil {
		return err
	}
	err = git.Tag(tagVersion, hash)
	if err != nil {
		return err
	}
	err = git.Push()
	if err != nil {
		return err
	}
	return nil
}
