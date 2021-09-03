package storage

import (
	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var linuxAmd64Download = types.Download{
	Os: "linux",
	Arch: "amd64",
}

var defaultConfig = config.Config{
	Base:  "./../../test_data/target",
	ArtifactDir: "./../../test_data/source",
	Namespace: "absaoss",
	TargetDir: "target",
	Branch: "gh-pages",
	WebRoot: "/",
}

func getDefaultPath() *path.Path{
	p, _ := path.NewPath(defaultConfig)
	return p
}

func TestCreatePlatformMetadata(t *testing.T) {
	// arrange

	fp, err := NewProvider(getDefaultPath())
	require.NoError(t, err)
	// act
	path, err := fp.CreatePlatformMetadata(linuxAmd64Download)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, defaultConfig.Base+"/absaoss/dummy/1.2.5/download/linux/amd64", path)
}

//todo: clean files
func TestMain(m *testing.M) {
	defer os.RemoveAll(defaultConfig.Base+"/absaoss")
	// before
	m.Run()
	// cleanup
}
