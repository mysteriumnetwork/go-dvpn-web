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

package command

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// PushChange pushes the change to github
func PushChange() error {
	gitToken := os.Getenv("GIT_TOKEN")
	if gitToken == "" {
		return errors.New("please specify the GIT_TOKEN environment variable")
	}

	tagVersion := os.Getenv("GIT_TAG_VERSION")
	if tagVersion == "" {
		return errors.New("please specify the TAG_VERSION environment variable")
	}

	fmt.Println("Committing...")
	repo, err := git.PlainOpen("./")
	if err != nil {
		return err
	}
	fmt.Println("repo opened")
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	fmt.Println("worktree fetched")

	fmt.Println("checking out master")
	branch := "refs/heads/master"
	b := plumbing.ReferenceName(branch)

	err = w.Checkout(&git.CheckoutOptions{
		Branch: b,
		Create: false,
		Force:  false,
	})
	if err != nil {
		return err
	}
	fmt.Println("master checked out")

	fmt.Println("adding changes")
	_, err = w.Add("assets_vfsdata.go")
	if err != nil {
		return err
	}
	fmt.Println("changes added")

	fmt.Println("performing commit")
	commitHash, err := w.Commit("Updating assets_vfsdata.go", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mister CI tool",
			Email: "dev@mysterium.network",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	fmt.Println("Commit done")

	fmt.Println("Tagging...", tagVersion)
	n := plumbing.ReferenceName("refs/tags/" + tagVersion)
	t := plumbing.NewHashReference(n, commitHash)
	err = repo.Storer.SetReference(t)
	if err != nil {
		return err
	}
	fmt.Println("tagged")

	fmt.Println("Pushing...")
	rs := config.RefSpec("refs/tags/*:refs/tags/*")
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			// this can be anything but not an empty string
			Username: "MisterFancyPants",
			Password: gitToken,
		},
		RefSpecs: []config.RefSpec{rs},
	})

	fmt.Println("Push done")

	return nil
}
