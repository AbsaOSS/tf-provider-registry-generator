package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	config2 "github.com/k0da/tfreg-golang/internal/config"
	pather "github.com/k0da/tfreg-golang/internal/dir"
	"github.com/k0da/tfreg-golang/internal/terraform"
	"github.com/absaoss/gopkg/shell"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func clone(c *config2.Config) {
	var cmd = shell.Command{
		Command: "git",
		Args:    []string{"clone", "--branch", c.Branch, c.RepoURL, "pages",
	}
	gitLog, err := shell.Execute(cmd)
	checkError(err)
	log.Infof("Git cloned with %s", gitLog)

	data, err := ioutil.ReadFile("data/terraform.json")
	checkError(err)
	dst := path + "/terraform.json"
	ioutil.WriteFile(dst, data, 0644)
	d1 := []byte("path: " + webroot)
	os.MkdirAll("pages/_data", 0755)
	os.WriteFile("pages/_data/root.yaml", d1, 0644)
}

func provider(c config.Config) {

	provider, err := terraform.NewProvider(c)
	checkError(err)
	err = provider.GenerateDownloadInfo()
	checkError(err)

	terraform.NewPath()
	f := terraform.NewFileProvider()

	pather, err := pather.NewDirProvider(config, provider)
	checkError(err)
	_, err = pather.CreateDownloadDirectory()
	checkError(err)
	versions, err := pather.GetVersions()
	checkError(err)
	version := provider.GenerateVersion()
	versions.Versions = append(versions.Versions, *version)
	pather.WriteVersions(versions)
	checkError(err)

}
func commit() {
	var lfsTrackCmd = shell.Command{
		Command: "git",
		Args:    []string{"lfs", "track", "download/*",
		Workdir: c.Base,
	}
	var gitAddAttrCmd = shell.Command{
		Command: "git",
		Args:    []string{"add", ".gitattributes",
		Workdir: c.Base,
	}
	gitAddLog, err := shell.Execute(gitAddAttrCmd)
	checkError(err)
	log.Infof("git add %s", gitAddLog)

	cmd = exec.Command("git", "commit", "-m", "auto")
	cmd.Stdout = &out
	err = cmd.Run()
	checkError(err)
	cmd = exec.Command("git", "push", "origin", os.Getenv("BRANCH"))
	cmd.Stdout = &out
	err = cmd.Run()
	checkError(err)
}

func main() {
	// clone pages repo
	config, err := config2.NewConfig("pages")
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
	//	fmt.Printf("Versions %+v", v.Versions)
	//	fmt.Println("\n")
	//	fmt.Printf("Download %+v", d)
}
