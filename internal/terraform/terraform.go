package terraform

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/k0da/tfreg-golang/internal/location"
	"github.com/k0da/tfreg-golang/internal/types"
)

const protocolVersion = "5.2"

type Provider struct {
	Platforms []types.Platform
	location  location.Locationer
}

func NewProvider(l location.Locationer) (p *Provider, err error) {
	if l == nil {
		err = fmt.Errorf("nil location provider")
		return
	}
	p = new(Provider)
	p.location = l
	for _, a := range p.location.GetArtifacts() {
		p.Platforms = append(p.Platforms, types.Platform{
			Os:         a.Os,
			Arch:       a.Arch,
			FileOrigin: a.File,
		})
	}
	return
}

func (p *Provider) GetDownloadInfo() (downloads []types.Download, err error) {
	// todo: testurl
	var url = p.location.UrlBinaries()
	var gpgPublicKey *types.GPGPublicKey
	downloads = []types.Download{}
	for _, platform := range p.Platforms {
		d := types.Download{Os: platform.Os, Arch: platform.Arch, Filename: platform.FileOrigin}
		d.DownloadURL = url + platform.FileOrigin
		d.Protocols = []string{protocolVersion}
		// todo: consider to check if files exists, don't necessarily to be on this place
		d.Shasum, err = getSHA256(p.location.ArtifactsPath() + "/" + platform.FileOrigin)
		if err != nil {
			return downloads, err
		}
		d.ShasumsSignatureURL = url + p.location.GetShaSumSignatureFile()
		d.ShasumsURL = url + p.location.GetShaSumFile()
		// todo: d.SigningKeys = resolve keys
		gpgPublicKey, err = p.ExtractPublicKey()
		if err != nil {
			return
		}
		d.SigningKeys.GpgPublicKeys = append(d.SigningKeys.GpgPublicKeys, *gpgPublicKey)
		downloads = append(downloads, d)
	}
	return
}

func (p *Provider) GenerateVersion() *types.Version {
	version := &types.Version{}
	version.Protocols = []string{protocolVersion}
	version.Version = p.location.GetVersion()
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
