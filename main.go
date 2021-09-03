package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/k0da/tfreg-golang/internal/storage"

	"github.com/k0da/tfreg-golang/internal/config"
	pather "github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/terraform"

	"github.com/AbsaOSS/gopkg/shell"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func clone(c config.Config) {
	var cmd = shell.Command{
		Command: "git",
		Args:    []string{"clone", "--branch", c.Branch, c.RepoURL, c.Base},
	}
	gitLog, err := shell.Execute(cmd)
	checkError(err)
	log.Printf("Git cloned with %s\n", gitLog)

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
	_ = provider.GenerateVersion()
	// ...
	//versions.Versions = append(versions.Versions, *version)
	//pather.WriteVersions(versions)
	checkError(err)

}
func commit(c config.Config) {
	var lfsTrackCmd = shell.Command{
		Command:    "git",
		Args:       []string{"lfs", "track", "download/*"},
		WorkingDir: c.Base,
	}
	lfsLog, err := shell.Execute(lfsTrackCmd)
	checkError(err)
	log.Printf("git lfs %s\n", lfsLog)
	var gitAddAttrCmd = shell.Command{
		Command:    "git",
		Args:       []string{"add", ".gitattributes"},
		WorkingDir: c.Base,
	}
	gitAddLog, err := shell.Execute(gitAddAttrCmd)
	checkError(err)
	log.Printf("git add %s\n", gitAddLog)
}

func main() {
	// clone pages repo
	config, err := config.NewConfig("pages")
	checkError(err)
	clone(config)
	provider(config)
	commit(config)
	//	ver, err := os.ReadFile("./version.json")
	//	if err != nil {
	//		fmt.Printf("Error %s", err.Error())
	//	}
	//	down, err := os.ReadFile("./download.json")
	//	if err != nil {
	//		fmt.Printf("Error %s", err.Error())
	//	}
	//	var v registry.Version
	//	var d registry.Download
	//	json.Unmarshal(ver, &v)
	//	err = json.Unmarshal(down, &d)
	//	if err != nil {
	//		fmt.Printf("Error %s", err.Error())
	//	}
	//	fmt.Printf("VersionsPath %+v", v.VersionsPath)
	//	fmt.Println("\n")
	//	fmt.Printf("Download %+v", d)
}
