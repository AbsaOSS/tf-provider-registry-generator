package terraform

import (
	"fmt"
	"github.com/AbsaOSS/gopkg/env"
	"github.com/k0da/tfreg-golang/config"
	"os"
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

type TerraformProvider struct {
	Name string
	Namespace string
	Version string
	Files []string
	Platforms []Platform
}

func (p *TerraformProvider) UpdatePlatform(os, arch string) *TerraformProvider {
	platform := Platform{Os: os, Arch: arch}
	p.Platforms = append(p.Platforms, platform)
	return p
}

func (p *TerraformProvider) GenerateArchs() {
	d := Download{}
	for _, platform := range p.Platforms {
		d.Os = platform.Os
		d.Arch = platform.Arch
	}
}

func(p *TerraformProvider) GenerateVersion() *Version {
	const protocolVersion = "5.2"
	version := &Version{}
	version.Protocols = []string{protocolVersion}
	version.Version = p.Version
	version.Platforms = p.Platforms
	return version
}


func NewProvider(c config.Config, name string) (p *TerraformProvider, err error) {
	var files []os.DirEntry
	files, err = os.ReadDir(c.ArtifactsPath())
	if err != nil {
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".zip") {
			file := f.Name()

			provider = parseProvider(file)
			provider.UpdatePlatform(getOs(file), getArch(file))
		}
	}

	p := new(TerraformProvider)
	// TODO: make validations of inputs here and initialize provide only with NewProvider
	p.Name = name
	return p
}

