package location

import (
	"fmt"
	"os"
	"strings"

	"github.com/AbsaOSS/tf-provider-registry-generator/internal/config"
)

type ILocation interface {
	ArtifactsPath() string
	TargetsPath() string
	VersionsPath() string
	DownloadsPath() string
	UrlBinaries() string
	GetArtifacts() []Artifact
	GetShaSumFile() string
	GetShaSumSignatureFile() string
	GetVersion() string
	TerraformJSONPath() string
	GetConfig() config.Config
}

type Location struct {
	artifacts []Artifact
	config    config.Config
	name      string
	version   string
}

type Artifact struct {
	Name    string
	Version string
	Os      string
	Arch    string
	File    string
}

func NewLocation(c config.Config) (l *Location, err error) {
	l = &Location{
		config: c,
	}
	var files []os.DirEntry
	files, err = os.ReadDir(c.ArtifactDir)
	if err != nil {
		return
	}
	l.artifacts, err = l.parseArtifacts(files)
	if err != nil {
		return
	}
	var m = map[string]bool{}
	for _, pi := range l.artifacts {
		m[pi.Name] = true
	}
	if len(m) != 1 {
		err = fmt.Errorf("more than one provider found in %s (%v)", l.config.ArtifactDir, m)
		return
	}
	l.name = l.artifacts[0].Name
	l.version = l.artifacts[0].Version
	return
}

func (l *Location) root() string {
	return l.config.Base
}

func (l *Location) providerRoot() string {
	return l.root() + "/" + l.config.Namespace + "/" + l.name
}

func (l *Location) ArtifactsPath() string {
	return l.config.ArtifactDir
}

func (l *Location) TargetsPath() string {
	return l.root() + "/" + l.config.TargetDir
}

func (l *Location) VersionsPath() string {
	return l.providerRoot() + "/versions"
}

func (l *Location) DownloadsPath() string {
	return l.providerRoot() + "/" + l.version + "/download"
}

func (l *Location) UrlBinaries() string {
	return "https://github.com/" + l.config.Owner + "/" + l.config.Repository + "/releases/download/" + l.GetVersion() + "/"
}

// GetArtifacts returns valid list of artifacts with at least one artifact
func (l *Location) GetArtifacts() []Artifact {
	return l.artifacts
}

func (l *Location) GetShaSumFile() string {
	return "terraform-provider-" + l.name + "_" + l.version + "_SHA256SUMS"
}

func (l *Location) GetShaSumSignatureFile() string {
	return "terraform-provider-" + l.name + "_" + l.version + "_SHA256SUMS.sig"
}

func (l *Location) GetVersion() string {
	return l.version
}

func (l *Location) TerraformJSONPath() string {
	return l.config.Base + "/terraform.json"
}

func (l *Location) GetConfig() config.Config {
	return l.config
}

// makes list of ArtifactsPath from files in the path
func (l *Location) parseArtifacts(files []os.DirEntry) (pis []Artifact, err error) {
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".zip") {
			var pi Artifact
			pi, err = l.getArtifactInfo(f.Name())
			if err != nil {
				return
			}
			pis = append(pis, pi)
		}
	}
	return
}

// parse name from file like this: terraform-provider-dummy_1.2.5_linux_amd64.zip into artifact
func (l *Location) getArtifactInfo(fileName string) (a Artifact, err error) {
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
	a.Name = strings.TrimPrefix(fileParts[0], "terraform-provider-")
	// todo: parse version (v1.2.3, 1.2.3 etc...), see: k8gb depresolver
	// todo: consider validations of OS and arch
	a.Version = fileParts[1]
	a.Os = fileParts[2]
	a.Arch = fileParts[3]
	a.File = fileName
	return
}
