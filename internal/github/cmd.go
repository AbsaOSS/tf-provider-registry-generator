package github

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/AbsaOSS/gopkg/shell"
	"github.com/k0da/tfreg-golang/internal/config"
)

const commitMsg = "Generate Terraform registry"

func Clone(c config.Config) (err error) {
	args := []string{"clone", "--branch", c.Branch, c.RepoURL, c.Base}
	err = runGitCmd(args, "")
	if err != nil {
		return
	}
	// to-do move to Location
	data, err := ioutil.ReadFile("data/terraform.json")
	if err != nil {
		return
	}
	dst := c.Base + "/terraform.json"
	ioutil.WriteFile(dst, data, 0644)
	d1 := []byte("path: " + c.WebRoot)
	os.MkdirAll(c.Base+"/_data", 0755)
	os.WriteFile(c.Base+"/_data/root.yaml", d1, 0644)

	return nil
}

func CommitAndPush(c config.Config) (err error) {
	lfsTrack := []string{"lfs", "track", "download/*"}
	err = runGitCmd(lfsTrack, c.Base)

	gitAddAttr := []string{"add", ".gitattributes"}
	err = runGitCmd(gitAddAttr, c.Base)

	gitUser := []string{"config", "user.name", c.User}
	err = runGitCmd(gitUser, c.Base)

	gitEmail := []string{"config", "user.email", c.Email}
	err = runGitCmd(gitEmail, c.Base)

	gitSetRemote := []string{"remote", "set-url", "origin", c.RepoURL}
	err = runGitCmd(gitSetRemote, c.Base)

	gitAdd := []string{"add", "./"}
	err = runGitCmd(gitAdd, c.Base)

	gitCommit := []string{"commit", "-m", commitMsg}
	err = runGitCmd(gitCommit, c.Base)

	gitPush := []string{"push", "origin", c.Branch}
	err = runGitCmd(gitPush, c.Base)

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
