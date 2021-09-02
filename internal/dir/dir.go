package dir

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/terraform"
)

// DirProvider retrieving proper path
type DirProvider struct {
	namespace string
	tp        *terraform.TerraformProvider
	path      *terraform.Path
}

func NewDirProvider(config config.Config, tp *terraform.TerraformProvider) (provider *DirProvider, err error) {
	if tp == nil {
		err = fmt.Errorf("nil terraform provider")
		return
	}
	provider = new(DirProvider)
	provider.tp = tp
	provider.path, err = terraform.NewPath(tp, config)
	return
}

func (p *DirProvider) CreateDownloadDirectory() (path string, err error) {
	err = os.MkdirAll(p.path.Downloads(), 0755)
	return
}

// GetVersions takes versions.json and retreives Versions struct
func (p *DirProvider) GetVersions() (v terraform.Versions, err error) {
	v = terraform.Versions{}
	data, err := os.ReadFile(p.path.Versions())
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

// WriteVersions stores versions.json
func (p *DirProvider) WriteVersions(v terraform.Versions) (err error) {
	if len(v.Versions) == 0 {
		err = fmt.Errorf("empty versions")
		return
	}
	data, err := json.Marshal(v)
	err = os.WriteFile(p.path.Versions(), data, 0644)
	return
}
