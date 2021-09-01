package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Download struct {
	Protocols           []string `json:"protocols"`
	Os                  string   `json:"os"`
	Arch                string   `json:"arch"`
	Filename            string   `json:"filename"`
	DownloadURL         string   `json:"download_url"`
	ShasumsURL          string   `json:"shasums_url"`
	ShasumsSignatureURL string   `json:"shasums_signature_url"`
	Shasum              string   `json:"shasum"`
	SigningKeys         struct {
		GpgPublicKeys []struct {
			KeyID          string `json:"key_id"`
			ASCIIArmor     string `json:"ascii_armor"`
			TrustSignature string `json:"trust_signature"`
			Source         string `json:"source"`
			SourceURL      string `json:"source_url"`
		} `json:"gpg_public_keys"`
	} `json:"signing_keys"`
}

type Versions struct {
	Versions []Version `json:"versions"`
}
type Version struct {
	Version   string     `json:"version"`
	Protocols []string   `json:"protocols"`
	Platforms []Platform `json:"platforms"`
}

type Platform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
}

type Provider struct {
	Name string
	Namespace string
	Version string
	Files []string
	Platforms []Platform
}

const TFProtocol = "5.2"

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

func getDownloadDir(p *Provider) string{
	return getNameDir(p.Name)+"/"+p.Version+"/download"
}

func parseProvider(file string) *Provider{
	prov := getProvider(file)
	p := new(Provider)
	p.Files = append(p.Files, file)
	p.Name = getProviderName(prov)
	p.Version = getVersion(prov)
	return p
}

func (p *Provider) updatePlatform(os, arch string) *Provider {
	platform := Platform{Os: os, Arch: arch}
	p.Platforms = append(p.Platforms, platform)
	return p
}

func (p *Provider) generateArchs() {
	d := Download{}
	for _, platform := range p.Platforms {
		d.Os = platform.Os
		d.Arch = platform.Arch
	}
}

func(p *Provider) generateVersion() *Version{
	version := &Version{}
	version.Protocols = []string{TFProtocol}
	version.Version = p.Version
	version.Platforms = p.Platforms
	return version
}

func prepareDownloadDir(p *Provider) error {
	// create pages/$NAMESPACE/$name/$version/download/
	return os.MkdirAll(getDownloadDir(p), 0755)
}

func newProvider(name string) *Provider {
	p := new(Provider)
	p.Name = name
	return p
}

func provider() {
	versions := Versions{}
	var provider *Provider
	files, err := os.ReadDir(path + "/" + os.Getenv("ARTIFACTS_DIR"))
	checkError(err)

	// walk through files and parse provider dist for name version os and arch
	for _, f := range files {
		if strings.Contains(f.Name(), ".zip") {
			file := f.Name()
			provider = parseProvider(file)
			provider.updatePlatform(getOs(file), getArch(file))
		}
	}
	provider.generateArchs()
	err = prepareDownloadDir(provider)
	checkError(err)
	version := provider.generateVersion()
	versionsPath := getVersionsPath(provider.Name)
	data, err := os.ReadFile(versionsPath)
	checkError(err)
	err = json.Unmarshal(data, &versions)
	checkError(err)
	// append existing versions
	versions.Versions = append(versions.Versions, *version)
	data, err = json.Marshal(versions)
	err = os.WriteFile(versionsPath, data, 0644)
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
