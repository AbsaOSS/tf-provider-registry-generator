package path

import (
	"encoding/json"
	"fmt"
	"github.com/AbsaOSS/gopkg/env"
	"github.com/k0da/tfreg-golang/terraform"
	"os"
)

// PathProvider retrieving proper path
type PathProvider struct {
	Base string
	namespace string
	tp *terraform.TerraformProvider
}

func NewPathProvider(base string, tp *terraform.TerraformProvider) (provider *PathProvider, err error) {
	const targetDir = "TARGET_DIR"
	if tp == nil {
		err = fmt.Errorf("nil terraform provider")
		return
	}
	target := env.GetEnvAsStringOrFallback(targetDir,"")
	if target == "" {
		err = fmt.Errorf("empty %s", targetDir)
		return
	}
	provider = new(PathProvider)
	provider.tp = tp
	provider.Base = fmt.Sprintf("%s/%s", base, target)
	return
}


func (p *PathProvider) CreateDownloadDirectory() (path string, err error){
	path = fmt.Sprintf("%s/%s/download",p.root(), p.tp.Version)
	err = os.MkdirAll(path, 0755)
	return
}

// GetVersions takes versions.json and retreives Versions struct
func (p *PathProvider) GetVersions() (v terraform.Versions, err error) {
	v = terraform.Versions{}
	path := fmt.Sprintf("%s/versions", p.root())
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

// WriteVersions stores versions.json
func (p *PathProvider) WriteVersions(v terraform.Versions) (err error) {
	if len(v.Versions) == 0 {
		err = fmt.Errorf("empty versions")
		return
	}
	path := fmt.Sprintf("%s/versions", p.root())
	data, err := json.Marshal(v)
	err = os.WriteFile(path, data, 0644)
	return
}


func (p *PathProvider) root() string {
	return fmt.Sprintf("%s/%s/%s",p.Base, p.tp.Namespace, p.tp.Name)
}

