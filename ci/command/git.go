package command

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// PushChange pushes the change to github
func PushChange() error {
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

	fmt.Println("adding changes")
	_, err = w.Add("assets_vfsdata.go")
	if err != nil {
		return err
	}
	fmt.Println("changes added")

	fmt.Println("performing commit")
	_, err = w.Commit("Updating assets_vfsdata.go", &git.CommitOptions{
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

	fmt.Println("Pushing...")
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "MisterFancyPants", // yes, this can be anything except an empty string
			Password: os.Getenv("GIT_TOKEN"),
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("Push done")

	return nil
}
