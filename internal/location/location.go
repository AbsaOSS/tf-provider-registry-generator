package location

import (
	"fmt"
	"os"
	"strings"

	"github.com/k0da/tfreg-golang/internal/config"
)

type Location struct {
	artifacts []Artifact
	config    config.Config
	Name      string
	Version   string
}

type Artifact struct {
	Name    string
	Version string
	Os      string
	Arch    string
	File    string
}

func NewLocation(c config.Config) (p *Location, err error) {
	p = &Location{
		config: c,
	}
	var files []os.DirEntry
	files, err = os.ReadDir(c.ArtifactDir)
	if err != nil {
		return
	}
	p.artifacts, err = p.parseArtifacts(files)
	if err != nil {
		return
	}
	var m = map[string]bool{}
	for _, pi := range p.artifacts {
		m[pi.Name] = true
	}
	if len(m) != 1 {
		err = fmt.Errorf("more than one provider found in %s (%v)", p.config.ArtifactDir, m)
		return
	}
	p.Name = p.artifacts[0].Name
	p.Version = p.artifacts[0].Version
	return
}

func (p *Location) root() string {
	return p.config.Base
}

func (p *Location) providerRoot() string {
	return p.root() + "/" + p.config.Namespace + "/" + p.Name
}

func (p *Location) ArtifactsPath() string {
	return p.config.ArtifactDir
}

func (p *Location) TargetsPath() string {
	return p.root() + "/" + p.config.TargetDir
}

func (p *Location) VersionsPath() string {
	return p.providerRoot() + "/versions"
}

func (p *Location) DownloadsPath() string {
	return p.providerRoot() + "/" + p.Version + "/download"
}

func (p *Location) BinariesPath() string {
	return p.root() + "/binaries"
}

func (p *Location) GPGPubring() string {
	return p.config.GPGHome + "/pubring.gpg"
}

func (p *Location) UrlBinaries() string {
	return "https://media.githubusercontent.com/media/" + p.config.Owner + "/" + p.config.Repository + "/" + p.config.Branch + "/binaries/"
}

// GetArtifacts returns valid list of artifacts with at least one artifact
func (p *Location) GetArtifacts() []Artifact {
	return p.artifacts
}

func (p *Location) GetShaSumFile() string {
	return "terraform-provider-" + p.Name + "_" + p.Version + "_SHA256SUMS"
}

func (p *Location) GetShaSumSignatureFile() string {
	return "terraform-provider-" + p.Name + "_" + p.Version + "_SHA256SUMS.sig"
}

func (p *Location) GPGFingerprint() string {
	return p.config.GPGFingerPrint
}


// makes list of ArtifactsPath from files in the path
func (p *Location) parseArtifacts(files []os.DirEntry) (pis []Artifact, err error) {
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".zip") {
			var pi Artifact
			pi, err = p.getArtifactInfo(f.Name())
			if err != nil {
				return
			}
			pis = append(pis, pi)
		}
	}
	return
}

// parse Name from file like this: terraform-provider-dummy_1.2.5_linux_amd64.zip into artifact
func (p *Location) getArtifactInfo(fileName string) (a Artifact, err error) {
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
