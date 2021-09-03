package storage

import (
	"encoding/json"
	"fmt"
	"github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/terraform"
	"os"
)

// FileProvider retrieving proper path
type FileProvider struct {
	namespace string
	path      *path.Path
}

type Filer interface {
	CreatePlatformMetadata(download terraform.Download) (path string, err error)
}

func NewFileProvider(path *path.Path) (provider *FileProvider, err error) {
	provider = new(FileProvider)
	provider.path = path
	return
}

func (p *FileProvider) CreatePlatformMetadata(download terraform.Download) (path string, err error) {
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
func (p *FileProvider) GetVersions() (v terraform.Versions, err error) {
	v = terraform.Versions{}
	data, err := os.ReadFile(p.path.VersionsPath())
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

// WriteVersions stores versions.json
func (p *FileProvider) WriteVersions(v terraform.Versions) (err error) {
	if len(v.Versions) == 0 {
		err = fmt.Errorf("empty versions")
		return
	}
	data, err := json.Marshal(v)
	err = os.WriteFile(p.path.VersionsPath(), data, 0644)
	return
}
