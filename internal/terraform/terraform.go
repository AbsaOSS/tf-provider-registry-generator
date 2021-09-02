package terraform

import (
	"fmt"
	"os"
	"strings"

	"github.com/k0da/tfreg-golang/internal/config"
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
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	fileName string
}

type TerraformProvider struct {
	Name      string
	Namespace string
	Version   string
	Platforms []Platform
}

type providerInfo struct {
	Name    string
	version string
	os      string
	arch    string
	file    string
}


func (p *TerraformProvider) GenerateDownloadInfo() (err error) {
	d := Download{}
	for _, platform := range p.Platforms {
		d.Os = platform.Os
		d.Arch = platform.Arch
	}
	// todo:
	return
}

func (p *TerraformProvider) GenerateVersion() *Version {
	const protocolVersion = "5.2"
	version := &Version{}
	version.Protocols = []string{protocolVersion}
	version.Version = p.Version
	for _, platform := range p.Platforms {
		version.Platforms = append(version.Platforms, Platform{Os: platform.Os, Arch: platform.Arch})
	}
	return version
}

func NewProvider(c config.Config) (p *TerraformProvider, err error) {
	var files []os.DirEntry
	files, err = os.ReadDir(c.ArtifactDir)
	if err != nil {
		return
	}
	pis, err := parseProviders(files)
	if err != nil {
		return
	}
	var m = map[string]bool{}
	for _, pi := range pis {
		m[pi.Name] = true
	}
	if len(m) != 1 {
		err = fmt.Errorf("more than one provider found in %s (%v)", c.ArtifactDir, m)
		return
	}
	p = new(TerraformProvider)
	p.Name = pis[0].Name
	p.Version = pis[0].version
	for _, pi := range pis {
		p.Platforms = append(p.Platforms, Platform{
			Os:       pi.os,
			Arch:     pi.arch,
			fileName: pi.file,
		})
	}
	return
}

// makes list of providerInfo from files in the path
func parseProviders(files []os.DirEntry) (pis []providerInfo, err error) {
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".zip") {
			var pi providerInfo
			pi, err = getProviderInfo(f.Name())
			if err != nil {
				return
			}
			pis = append(pis, pi)
		}
	}
	return
}

// parse name from file like this: terraform-provider-dummy_1.2.5_linux_amd64.zip
func getProviderInfo(fileName string) (i providerInfo, err error) {
	const prefix = "terraform-provider-"
	if !strings.HasPrefix(fileName, prefix) {
		err = fmt.Errorf("filed to parse %s, must start with %s", fileName, prefix)
		return
	}
	t := strings.TrimSuffix(fileName, ".zip")
	fileParts := strings.Split(t, "_")
	if len(fileParts) != 4 {
		err = fmt.Errorf("filed to parse %s, expecting %s_<version>_<os>_<arch>.zip", fileName, prefix)
		return
	}
	i.Name = strings.TrimPrefix(fileParts[0], "terraform-provider-")
	// todo: parse version (v1.2.3, 1.2.3 etc...), see: k8gb depresolver
	// todo: consider validations of OS and arch
	i.version = fileParts[1]
	i.os = fileParts[2]
	i.arch = fileParts[3]
	i.file = fileName
	return
}
