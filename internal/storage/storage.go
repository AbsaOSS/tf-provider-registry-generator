package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/types"
)

// Provider retrieving proper path
type Provider struct {
	namespace string
	path      *path.Path
}

type Storage interface {
	CreatePlatformMetadata(download types.Download) (path string, err error)
	GetVersions() (v types.Versions, err error)
	WriteVersions(v types.Versions) (err error)
}

func NewProvider(path *path.Path) (provider *Provider, err error) {
	provider = new(Provider)
	provider.path = path
	return
}

func (p *Provider) CreatePlatformMetadata(download types.Download) (path string, err error) {
	dir := p.path.DownloadsPath() + "/" + download.Os
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	path = dir + "/" + download.Arch
	b, err := json.Marshal(download)
	if err != nil {
		return
	}
	err = os.WriteFile(path, b, 0644)
	return
}

// GetVersions takes versions.json and retreives Versions struct
// if file doesn't exists, return empty Versions slice
func (p *Provider) GetVersions() (v types.Versions, err error) {
	v = types.Versions{}
	if _, err := os.Stat(p.path.VersionsPath()); os.IsNotExist(err) {
		return v, nil
	}
	data, err := os.ReadFile(p.path.VersionsPath())
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
	err = os.WriteFile(p.path.VersionsPath(), data, 0644)
	return
}

func (p *Provider) SaveBinaries() (err error) {
	err = os.MkdirAll(p.path.BinariesPath(), 0755)
	if err != nil {
		return
	}
	for _, a := range p.path.GetArtifacts() {
		err = p.copy(a.File)
		if err != nil {
			return err
		}
	}
	err = p.copy(p.path.GetShaSumFile())
	if err != nil {
		return err
	}
	err = p.copy(p.path.GetShaSumSignatureFile())
	return
}

func (p *Provider) copy(file string) (err error) {
	src := p.path.ArtifactsPath() + "/" + file
	dst := p.path.BinariesPath() + "/" + file
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
