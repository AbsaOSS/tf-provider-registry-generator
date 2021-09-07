package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/k0da/tfreg-golang/internal/cmd"
	"github.com/k0da/tfreg-golang/internal/config"
	pather "github.com/k0da/tfreg-golang/internal/location"
	"github.com/k0da/tfreg-golang/internal/storage"
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
	err := cmd.Run("git", args, "")
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
	location, err := pather.NewLocation(c)
	checkError(err)
	storage, err := storage.NewProvider(location)
	checkError(err)
	provider, err := terraform.NewProvider(location)
	checkError(err)
	downloads, err := provider.GetDownloadInfo()
	checkError(err)
	_, err = storage.WritePlatformMetadata(downloads)
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
	err := cmd.Run("git", lfsTrack, c.Base)
	checkError(err)
	gitAddAttr := []string{"add", ".gitattributes"}
	err = cmd.Run("git", gitAddAttr, c.Base)
	checkError(err)
	gitUser := []string{"config", "user.name", c.User}
	err = cmd.Run("git", gitUser, c.Base)
	checkError(err)
	gitEmail := []string{"config", "user.email", c.Email}
	err = cmd.Run("git", gitEmail, c.Base)
	checkError(err)
	gitSetRemote := []string{"remote" ,"set-url", "origin", c.RepoURL}
	err = cmd.Run("git", gitSetRemote, c.Base)
	checkError(err)
	gitAdd := []string{"add", "./"}
	err = cmd.Run("git", gitAdd, c.Base)
	checkError(err)
	gitCommit := []string{"commit", "-m", commitMsg}
	err = cmd.Run("git",gitCommit, c.Base)
	checkError(err)
	gitPush := []string{"push", "origin", c.Branch}
	err = cmd.Run("git", gitPush, c.Base)
	checkError(err)
}

func main() {
	config, err := config.NewConfig("pages")
	checkError(err)
	clone(config)
	provider(config)
	commit(config)
}
