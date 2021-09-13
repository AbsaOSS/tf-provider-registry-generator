package repo

import (
	"fmt"
	"log"

	"github.com/AbsaOSS/tf-provider-registry-generator/internal/location"

	"github.com/AbsaOSS/gopkg/shell"
)

const commitMsg = "Generate Terraform registry"

type IRepo interface {
	Clone() (err error)
	CommitAndPush() (err error)
}

type Github struct {
	location location.ILocation
}

func NewGithub(location location.ILocation) (g *Github, err error) {
	//todo: remove config
	if location == nil {
		return nil, fmt.Errorf("nil location")
	}
	g = &Github{
		location: location,
	}
	return
}

func (g *Github) Clone() (err error) {
	// todo: use location instead of config
	args := []string{"clone", "--branch", g.location.GetConfig().Branch, g.location.GetConfig().RepoURL, g.location.GetConfig().Base}
	err = runGitCmd(args, "")
	return
}

func (g *Github) CommitAndPush() (err error) {
	gitUser := []string{"config", "user.name", g.location.GetConfig().User}
	err = runGitCmd(gitUser, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitEmail := []string{"config", "user.email", g.location.GetConfig().Email}
	err = runGitCmd(gitEmail, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitSetRemote := []string{"remote", "set-url", "origin", g.location.GetConfig().RepoURL}
	err = runGitCmd(gitSetRemote, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitAdd := []string{"add", "./"}
	err = runGitCmd(gitAdd, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitCommit := []string{"commit", "-m", commitMsg}
	err = runGitCmd(gitCommit, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitPush := []string{"push", "origin", g.location.GetConfig().Branch}
	err = runGitCmd(gitPush, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	return
}

func runGitCmd(args []string, dir string) (err error) {
	command := shell.Command{
		Command:    "git",
		Args:       args,
		WorkingDir: dir,
	}
	cmdLog, err := shell.Execute(command)
	if err != nil {
		return
	}
	log.Printf("run git %+v: %s", args, cmdLog)
	return
}
