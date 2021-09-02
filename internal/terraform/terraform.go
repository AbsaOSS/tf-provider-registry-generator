package terraform

import (
	"crypto/sha256"
	"fmt"
	"github.com/k0da/tfreg-golang/internal/config"
	"os"
	"strings"
)
const protocolVersion = "5.2"

type TerraformProvider struct {
	name      string
	Version   string
	Platforms []Platform
	path         *Path
	fileProvider Filer
}


type providerInfo struct {
	Name    string
	version string
	os      string
	arch    string
	file    string
}

func NewProvider(c config.Config, fp Filer) (p *TerraformProvider, err error) {
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
	p.path, err = NewPath(p,c)
	if err != nil {
		return
	}
	p.name = pis[0].Name
	p.Version = pis[0].version
	for _, pi := range pis {
		p.Platforms = append(p.Platforms, Platform{
			Os:       pi.os,
			Arch:     pi.arch,
			fileName: pi.file,
		})
	}
	if fp == nil {
		err = fmt.Errorf("nil file provider")
		return
	}
	p.fileProvider = fp
	return
}

func (p *TerraformProvider) GenerateDownloadInfo() (err error) {
	const url = "https://media.githubusercontent.com/media/downloads/"
	var path string
	for _, platform := range p.Platforms {
		d := Download{Os: platform.Os, Arch: platform.Arch, Filename: platform.fileName}
		d.DownloadURL = url+platform.fileName
		d.Protocols = []string{protocolVersion}
		d.Shasum, err = getSHA256(p.path.Artifacts() + "/" + platform.fileName)
		if err != nil {
			return err
		}
		d.ShasumsSignatureURL = url + "terraform-provider-"+p.name+"_"+p.Version+"_SHA256SUMS.sig"
		d.ShasumsURL = url + "terraform-provider-"+p.name+"_"+p.Version+"_SHA256SUMS"
		// todo: d.SigningKeys = resolve keys
		path, err = p.fileProvider.CreatePlatformMetadata(d)
		if err != nil {
			err = fmt.Errorf("error writing metadata %s, %s", path, err)
			return
		}
	}
	return
}

func (p *TerraformProvider) GenerateVersion() *Version {
	version := &Version{}
	version.Protocols = []string{protocolVersion}
	version.Version = p.Version
	for _, platform := range p.Platforms {
		version.Platforms = append(version.Platforms, Platform{Os: platform.Os, Arch: platform.Arch})
	}
	return version
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


func getSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:]), err
}
