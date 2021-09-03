package terraform

import (
	"crypto/sha256"
	"fmt"
	"github.com/k0da/tfreg-golang/internal/path"
	"os"
)
const protocolVersion = "5.2"

type Provider struct {
	Platforms []Platform
	path       *path.Path
	fileProvider Filer
}

func NewProvider(path *path.Path, fp Filer) (p *Provider, err error) {
	if path == nil {
		err = fmt.Errorf("nil path provider")
		return
	}
	if fp == nil {
		err = fmt.Errorf("nil file provider")
		return
	}
	p = new(Provider)
	p.fileProvider = fp
	p.path = path
	for _, a := range p.path.GetArtifacts() {
		p.Platforms = append(p.Platforms, Platform{
			Os:       a.Os,
			Arch:     a.Arch,
			fileName: a.File,
		})
	}
	return
}

func (p *Provider) GenerateDownloadInfo() (err error) {
	const url = "https://media.githubusercontent.com/media/downloads/"
	var path string
	for _, platform := range p.Platforms {
		d := Download{Os: platform.Os, Arch: platform.Arch, Filename: platform.fileName}
		d.DownloadURL = url+platform.fileName
		d.Protocols = []string{protocolVersion}
		d.Shasum, err = getSHA256(p.path.ArtifactsPath() + "/" + platform.fileName)
		if err != nil {
			return err
		}
		d.ShasumsSignatureURL = url + "terraform-provider-"+p.path.Name+"_"+p.path.Version+"_SHA256SUMS.sig"
		d.ShasumsURL = url + "terraform-provider-"+p.path.Name+"_"+p.path.Version+"_SHA256SUMS"
		// todo: d.SigningKeys = resolve keys
		path, err = p.fileProvider.CreatePlatformMetadata(d)
		if err != nil {
			err = fmt.Errorf("error writing metadata %s, %s", path, err)
			return
		}
	}
	return
}

func (p *Provider) GenerateVersion() *Version {
	version := &Version{}
	version.Protocols = []string{protocolVersion}
	version.Version = p.path.Version
	for _, platform := range p.Platforms {
		version.Platforms = append(version.Platforms, Platform{Os: platform.Os, Arch: platform.Arch})
	}
	return version
}

func getSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:]), err
}
