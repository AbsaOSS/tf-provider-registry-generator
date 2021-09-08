package repo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/k0da/tfreg-golang/internal/location"

	"github.com/AbsaOSS/gopkg/shell"
	"github.com/k0da/tfreg-golang/internal/config"
)

const commitMsg = "Generate Terraform registry"

type IRepo interface {
	Clone() (err error)
	CommitAndPush() (err error)
}

type Github struct {
	location location.ILocation
	config config.Config
}

func NewGithub(c config.Config, location location.ILocation) (g *Github, err error) {
	//todo: remove config
	if location == nil {
		return nil, fmt.Errorf("nil location")
	}
	g = &Github{
		location: location,
		config: c,
	}
	return
}

func (g *Github) Clone() (err error) {
	// todo: use location instead of config
	args := []string{"clone", "--branch", g.config.Branch, g.config.RepoURL, g.config.Base}
	err = runGitCmd(args, "")
	if err != nil {
		return
	}
	// to-do move to Location
	data, err := ioutil.ReadFile("data/terraform.json")
	if err != nil {
		return
	}
	dst := g.config.Base + "/terraform.json"
	ioutil.WriteFile(dst, data, 0644)
	d1 := []byte("path: " + g.config.WebRoot)
	os.MkdirAll(g.config.Base+"/_data", 0755)
	os.WriteFile(g.config.Base+"/_data/root.yaml", d1, 0644)

	return nil
}

func (g *Github) CommitAndPush() (err error) {
	lfsTrack := []string{"lfs", "track", "download/*"}
	err = runGitCmd(lfsTrack, g.config.Base)

	gitAddAttr := []string{"add", ".gitattributes"}
	err = runGitCmd(gitAddAttr, g.config.Base)

	gitUser := []string{"config", "user.name", g.config.User}
	err = runGitCmd(gitUser, g.config.Base)

	gitEmail := []string{"config", "user.email", g.config.Email}
	err = runGitCmd(gitEmail, g.config.Base)

	gitSetRemote := []string{"remote", "set-url", "origin", g.config.RepoURL}
	err = runGitCmd(gitSetRemote, g.config.Base)

	gitAdd := []string{"add", "./"}
	err = runGitCmd(gitAdd, g.config.Base)

	gitCommit := []string{"commit", "-m", commitMsg}
	err = runGitCmd(gitCommit, g.config.Base)

	gitPush := []string{"push", "origin",g.config.Branch}
	err = runGitCmd(gitPush, g.config.Base)

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
