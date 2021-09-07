package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/k0da/tfreg-golang/internal/location"
	"github.com/k0da/tfreg-golang/internal/types"
)

// Provider retrieving proper location
type Provider struct {
	namespace string
	location  location.ILocation
}

type IStorage interface {
	WritePlatformMetadata(download types.Download) (path string, err error)
	GetVersions() (v types.Versions, err error)
	WriteVersions(v types.Versions) (err error)
}

func NewProvider(l location.ILocation) (provider *Provider, err error) {
	provider = new(Provider)
	provider.location = l
	return
}

func (p *Provider) WritePlatformMetadata(downloads []types.Download) (path string, err error) {
	var b []byte
	for _, d := range downloads {
		dir := p.location.DownloadsPath() + "/" + d.Os
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return
		}
		path = dir + "/" + d.Arch
		b, err = json.Marshal(d)
		if err != nil {
			return
		}
		err = os.WriteFile(path, b, 0644)
		if err != nil {
			return
		}
	}
	return
}

// GetVersions takes versions.json and retreives Versions struct
// if file doesn't exists, return empty Versions slice
func (p *Provider) GetVersions() (v types.Versions, err error) {
	v = types.Versions{}
	if _, err := os.Stat(p.location.VersionsPath()); os.IsNotExist(err) {
		return v, nil
	}
	data, err := os.ReadFile(p.location.VersionsPath())
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

// WriteVersions stores versions.json
func (p *Provider) WriteVersions(v types.Versions) (err error) {
	if len(v.Versions) == 0 {
		err = fmt.Errorf("empty versions")
		return
	}
	data, err := json.Marshal(v)
	err = os.WriteFile(p.location.VersionsPath(), data, 0644)
	return
}

func (p *Provider) SaveBinaries() (err error) {
	err = os.MkdirAll(p.location.BinariesPath(), 0755)
	if err != nil {
		return
	}
	for _, a := range p.location.GetArtifacts() {
		err = p.copy(a.File)
		if err != nil {
			return err
		}
	}
	err = p.copy(p.location.GetShaSumFile())
	if err != nil {
		return err
	}
	err = p.copy(p.location.GetShaSumSignatureFile())
	return
}

func (p *Provider) copy(file string) (err error) {
	src := p.location.ArtifactsPath() + "/" + file
	dst := p.location.BinariesPath() + "/" + file
	log.Printf("copying file from %s to %s", src, dst)
	_, err = Copy(src, dst)
	return err
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
