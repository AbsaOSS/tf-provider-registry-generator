package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/k0da/tfreg-golang/internal/storage"

	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/git"
	pather "github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/terraform"

)

const commitMsg = "Generate Terraform registry"

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func clone(c config.Config) {
	args := []string{"clone", "--branch", c.Branch, c.RepoURL, c.Base}
	err := git.RunGit(args, "")
	checkError(err)

	data, err := ioutil.ReadFile("data/terraform.json")
	checkError(err)
	dst := c.Base + "/terraform.json"
	ioutil.WriteFile(dst, data, 0644)
	d1 := []byte("path: " + c.WebRoot)
	os.MkdirAll(c.Base+"/_data", 0755)
	os.WriteFile(c.Base+"/_data/root.yaml", d1, 0644)
}

func provider(c config.Config) {
	// init
	path, err := pather.NewPath(c)
	checkError(err)
	storage, err := storage.NewProvider(path)
	checkError(err)
	provider, err := terraform.NewProvider(path, storage)
	checkError(err)
	err = provider.GenerateDownloadInfo()
	checkError(err)
	versions, err := storage.GetVersions()
	checkError(err)
	version := provider.GenerateVersion()
	versions.Versions = append(versions.Versions, *version)
	err = storage.WriteVersions(versions)
	checkError(err)
	err = storage.SaveBinaries()
	checkError(err)
}

func commit(c config.Config) {
	lfsTrack := []string{"lfs", "track", "download/*"}
	err := git.RunGit(lfsTrack, c.Base)
	checkError(err)
	gitAddAttr := []string{"add", ".gitattributes"}
	err = git.RunGit(gitAddAttr, c.Base)
	checkError(err)
	gitUser := []string{"config", "user.name"}
	err = git.RunGit(gitUser, c.Base)
	checkError(err)
	gitEmail := []string{"config", "user.email"}
	err = git.RunGit(gitEmail, c.Base)
	checkError(err)
	gitSetRemote := []string{"remote" ,"set-url", "origin", c.RepoURL}
	err = git.RunGit(gitSetRemote, c.Base)
	checkError(err)
	gitAdd := []string{"add", "./"}
	err = git.RunGit(gitAdd, c.Base)
	checkError(err)
	gitCommit := []string{"commit", "-m", commitMsg}
	err = git.RunGit(gitCommit, c.Base)
	checkError(err)
	gitPush := []string{"push", "origin", c.Branch}
	err = git.RunGit(gitPush, c.Base)
	checkError(err)
}

func main() {
	config, err := config.NewConfig("pages")
	checkError(err)
	clone(config)
	provider(config)
	commit(config)
}
