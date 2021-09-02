package terraform

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileProvider retrieving proper path
type FileProvider struct {
	namespace string
	path      *Path
}

type Filer interface {
	CreatePlatformMetadata(download Download) (path string, err error)
}

func NewFileProvider(path *Path) (provider *FileProvider, err error) {
	provider = new(FileProvider)
	provider.path = path
	return
}

func (p *FileProvider) CreatePlatformMetadata(download Download) (path string, err error) {
	dir := p.path.Downloads()+"/"+download.Os
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
func (p *FileProvider) GetVersions() (v Versions, err error) {
	v = Versions{}
	data, err := os.ReadFile(p.path.Versions())
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

// WriteVersions stores versions.json
func (p *FileProvider) WriteVersions(v Versions) (err error) {
	if len(v.Versions) == 0 {
		err = fmt.Errorf("empty versions")
		return
	}
	data, err := json.Marshal(v)
	err = os.WriteFile(p.path.Versions(), data, 0644)
	return
}
