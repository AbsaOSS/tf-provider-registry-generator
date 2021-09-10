package terraform

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/AbsaOSS/tf-provider-registry-generator/internal/location"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/types"
)

const protocolVersion = "5.2"

type ITerraform interface {
	GetDownloadInfo() (downloads []types.Download, err error)
	GenerateVersion() *types.Version
	GenerateTerraformJSON() error
}

type TerraformProvider struct {
	Platforms []types.Platform
	location  location.ILocation
}

func NewTerraformProvider(l location.ILocation) (p *TerraformProvider, err error) {
	if l == nil {
		err = fmt.Errorf("nil location provider")
		return
	}
	p = new(TerraformProvider)
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

func (p *TerraformProvider) GetDownloadInfo() (downloads []types.Download, err error) {
	// todo: testurl
	var url = p.location.UrlBinaries()
	var keyUnquote string
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
		keyUnquote, err = strconv.Unquote(`"` +p.location.GetConfig().GPGArmor+`"`)
		if err != nil {
			return downloads, err
		}
		gpgPublicKey := types.GPGPublicKey{
			KeyID:      p.location.GetConfig().GPGKeyID,
			ASCIIArmor: keyUnquote,
		}

		d.SigningKeys.GpgPublicKeys = append(d.SigningKeys.GpgPublicKeys, gpgPublicKey)
		downloads = append(downloads, d)
	}
	return
}

func (p *TerraformProvider) GenerateVersion() *types.Version {
	version := &types.Version{}
	version.Protocols = []string{protocolVersion}
	version.Version = p.location.GetVersion()
	for _, platform := range p.Platforms {
		version.Platforms = append(version.Platforms, types.Platform{Os: platform.Os, Arch: platform.Arch})
	}
	return version
}

func (p *TerraformProvider) GenerateTerraformJSON() (err error) {
	const dataPrefix = "---\npermalink: /.well-known/terraform.json\n---\n"
	var data = fmt.Sprintf("%s{\"providers.v1\":\"%s\"}\n", dataPrefix, p.location.GetConfig().WebRoot)
	err = ioutil.WriteFile(p.location.TerraformJSONPath(), []byte(data), 0644)
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
