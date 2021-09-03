package terraform

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/storage"
	"github.com/k0da/tfreg-golang/internal/types"
)

const protocolVersion = "5.2"

type Provider struct {
	Platforms    []types.Platform
	path         *path.Path
	fileProvider storage.Storage
}

func NewProvider(path *path.Path, fp storage.Storage) (p *Provider, err error) {
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
		p.Platforms = append(p.Platforms, types.Platform{
			Os:         a.Os,
			Arch:       a.Arch,
			FileOrigin: a.File,
		})
	}
	return
}

func (p *Provider) GenerateDownloadInfo() (err error) {
	// todo: testurl
	var url = p.path.UrlBinaries()
	var path string
	for _, platform := range p.Platforms {
		d := types.Download{Os: platform.Os, Arch: platform.Arch, Filename: platform.FileOrigin}
		d.DownloadURL = url + platform.FileOrigin
		d.Protocols = []string{protocolVersion}
		// todo: consider to check if files exists, don't necessarily to be on this place
		d.Shasum, err = getSHA256(p.path.ArtifactsPath() + "/" + platform.FileOrigin)
		if err != nil {
			return err
		}
		d.ShasumsSignatureURL = url + p.path.GetShaSumSignatureFile()
		d.ShasumsURL = url + p.path.GetShaSumFile()
		// todo: d.SigningKeys = resolve keys
		path, err = p.fileProvider.CreatePlatformMetadata(d)
		if err != nil {
			err = fmt.Errorf("error writing metadata %s, %s", path, err)
			return
		}
	}
	return
}

func (p *Provider) GenerateVersion() *types.Version {
	version := &types.Version{}
	version.Protocols = []string{protocolVersion}
	version.Version = p.path.Version
	for _, platform := range p.Platforms {
		version.Platforms = append(version.Platforms, types.Platform{Os: platform.Os, Arch: platform.Arch})
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
