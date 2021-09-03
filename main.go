package main

import (
	"bytes"
	"github.com/k0da/tfreg-golang/internal/storage"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/k0da/tfreg-golang/internal/config"
	pather "github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/terraform"
)

var path = "pages/" + os.Getenv("TARGET_DIR")

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func clone() {
	//path := "pages/" + os.Getenv("TARGET_DIR")
	webroot := "path: " + os.Getenv("WEB_ROOT")
	cmd := exec.Command("git", "clone", "--branch", os.Getenv("BRANCH"), os.Getenv("REPO_URL"), "pages")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	checkError(err)

	data, err := ioutil.ReadFile("/data/terraform.json")
	checkError(err)
	dst := path + "/terraform.json"
	ioutil.WriteFile(dst, data, 0644)
	d1 := []byte(webroot)
	os.MkdirAll(path+"/_data", 0755)
	os.WriteFile(path+"/_data/root.yaml", d1, 0644)
}

func provider() {

	// init
	config, err := config.NewConfig("pages")
	checkError(err)
	path, err := pather.NewPath(config)
	checkError(err)
	file, err := storage.NewProvider(path)
	checkError(err)
	provider, err := terraform.NewProvider(path,file)
	checkError(err)
	err = provider.GenerateDownloadInfo()
	checkError(err)

	_ = provider.GenerateVersion()
	//versions.Versions = append(versions.Versions, *version)
	//pather.WriteVersions(versions)
	checkError(err)

}
func commit() {
	path := "pages/" + os.Getenv("TARGET_DIR")
	cmd := exec.Command("git", "add", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	checkError(err)
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
	clone()
	provider()
	commit()
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
