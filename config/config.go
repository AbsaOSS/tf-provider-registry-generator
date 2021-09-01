package config

import (
	"fmt"
	"github.com/AbsaOSS/gopkg/env"
)

type Config struct {
	TargetDir string
	Namespace string
	ArtifactDir string
	Branch string
	WebRoot string
	Base string
}

func NewConfig(base string) (c Config, err error) {
	const targetDir = "TARGET_DIR"
	const artifactsDir = "ARTIFACTS_DIR"
	const namespace = "NAMESPACE"
	c = Config{}
	c.TargetDir = env.GetEnvAsStringOrFallback(targetDir,"")
	if c.TargetDir == "" {
		err = fmt.Errorf("empty %s", targetDir)
		return
	}
	c.ArtifactDir = env.GetEnvAsStringOrFallback(artifactsDir,"")
	if c.ArtifactDir == "" {
		err = fmt.Errorf("empty %s", c.ArtifactDir)
		return
	}
	c.Namespace = env.GetEnvAsStringOrFallback(namespace,"")
	if c.Namespace == "" {
		err = fmt.Errorf("empty %s", namespace)
		return
	}
	if base == "" {
		err = fmt.Errorf("empty base")
		return
	}
	c.Base = base
	return
}

// todo: create class which read config and returns bunch of paths
func (c Config) BasePath() string {
	return c.Base
}

func (c Config) ArtifactsPath() string {
	return c.BasePath() + "/" + c.ArtifactDir
}

