package main

import (
	"bytes"
	pather "github.com/k0da/tfreg-golang/path"
	"github.com/k0da/tfreg-golang/terraform"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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

func getProvider(file string) string {
	return strings.TrimSuffix(file, ".zip")
}

func getOs(file string) string {
        name := getProvider(file)
	fileParts := strings.Split(name, "_")
	return fileParts[2]
}
func getArch(file string) string {
        name := getProvider(file)
	fileParts := strings.Split(name, "_")
	return fileParts[3]
}
func getVersion(name string) string {
	fileParts := strings.Split(name, "_")
	return fileParts[1]
}
func getProviderName(name string) string {
	fileParts := strings.Split(name, "_")
	return strings.TrimPrefix(fileParts[0], "terraform-provider-")
}

func getVersionsPath(name string) string {
	return getNameDir(name) + "/versions"
}

func getNameDir(name string) string {
	return path + "/" + os.Getenv("NAMESPACE") + "/" + name
}

func parseProvider(file string) *terraform.TerraformProvider{
	prov := getProvider(file)
	p := new(terraform.TerraformProvider)
	p.Files = append(p.Files, file)
	p.Name = getProviderName(prov)
	p.Version = getVersion(prov)
	return p
}

func provider() {
	var provider *terraform.TerraformProvider
	files, err := os.ReadDir(path + "/" + os.Getenv("ARTIFACTS_DIR"))
	checkError(err)

	// walk through files and parse provider dist for name version os and arch
	for _, f := range files {
		if strings.Contains(f.Name(), ".zip") {
			file := f.Name()
			provider = parseProvider(file)
			provider.UpdatePlatform(getOs(file), getArch(file))
		}
	}
	provider.GenerateArchs()

	pather,err := pather.NewPathProvider("pages", provider)
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
	//	fmt.Printf("Versions %+v", v.Versions)
	//	fmt.Println("\n")
	//	fmt.Printf("Download %+v", d)
}
