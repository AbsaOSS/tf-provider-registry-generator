package storage

import (
	"github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var linuxAmd64Download = terraform.Download{
	Os: "linux",
	Arch: "amd64",
}

func getDefaultPath() *path.Path{
	p, _ := path.NewPath(defaultConfig)
	return p
}

func TestCreatePlatformMetadata(t *testing.T) {
	// arrange

	fp, err := NewFileProvider(getDefaultPath())
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
