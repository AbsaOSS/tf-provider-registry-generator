package terraform

import (
	"fmt"

	"github.com/k0da/tfreg-golang/internal/config"
)

type Path struct {
	config   config.Config
	provider *TerraformProvider
}

func NewPath(provider *TerraformProvider, config config.Config) (p *Path, err error) {
	if provider == nil {
		err = fmt.Errorf("nil terraform provider")
		return
	}
	p = &Path{
		config:   config,
		provider: provider,
	}
	return
}

// todo: create class which read config and returns bunch of paths
func (p *Path) Root() string {
	return p.config.Base
}

func (p *Path) providerRoot() string {
	return p.Root() + "/" + p.provider.Namespace + "/" + p.provider.Name
}

func (p *Path) Artifacts() string {
	return p.Root() + "/" + p.config.ArtifactDir
}

func (p *Path) Targets() string {
	return p.Root() + "/" + p.config.TargetDir
}

func (p *Path) Versions() string {
	return p.providerRoot() + "/versions"
}

func (p *Path) Downloads() string {
	return p.providerRoot() + "/" + p.provider.Version + "/download"
}
