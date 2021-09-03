package storage

import (
	"encoding/json"
	"fmt"
	"github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/types"
	"os"
)

// Provider retrieving proper path
type Provider struct {
	namespace string
	path      *path.Path
}

type Storage interface {
	CreatePlatformMetadata(download types.Download) (path string, err error)
}

func NewProvider(path *path.Path) (provider *Provider, err error) {
	provider = new(Provider)
	provider.path = path
	return
}

func (p *Provider) CreatePlatformMetadata(download types.Download) (path string, err error) {
	dir := p.path.DownloadsPath()+"/"+download.Os
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	path = dir + "/" + download.Arch
	b,err := json.Marshal(download)
	if err != nil {
		return
	}
	err = os.WriteFile(path, b, 0644)
	return
}

// GetVersions takes versions.json and retreives Versions struct
func (p *Provider) GetVersions() (v types.Versions, err error) {
	v = types.Versions{}
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
