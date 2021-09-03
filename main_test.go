package main

import (
	"testing"

	"github.com/k0da/tfreg-golang/internal/config"
	pather "github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/storage"
	"github.com/k0da/tfreg-golang/internal/terraform"
	"github.com/stretchr/testify/require"
)

//func TestGenerateDownloadPath(t *testing.T) {
//	path, err := os.MkdirTemp("./", "test-pages-")
//	assert.NoError(t, err)
//	err = os.Setenv("NAMESPACE", "absaoss")
//	assert.NoError(t, err)
//	prepareDownloadDir(expectedProvider)
//	t.Logf("got path %s", path)
//}

var defaultConfig = config.Config{
	// todo: find a system to store / restore files in test
	Base:        "./test_data/target_green",
	ArtifactDir: "./test_data/source",
	Namespace:   "absaoss",
	Branch:      "gh-pages",
	WebRoot:     "/",
	Owner:       "absaoss",
	Repository:  "terraform-provider-dummy",
}

func TestGreenPath(t *testing.T) {
	// init
	path, err := pather.NewPath(defaultConfig)
	require.NoError(t, err)
	storage, err := storage.NewProvider(path)
	require.NoError(t, err)
	provider, err := terraform.NewProvider(path, storage)
	require.NoError(t, err)
	err = provider.GenerateDownloadInfo()
	require.NoError(t, err)
	versions, err := storage.GetVersions()
	require.NoError(t, err)
	version := provider.GenerateVersion()
	versions.Versions = append(versions.Versions, *version)
	err = storage.WriteVersions(versions)
	require.NoError(t, err)
	err = storage.SaveBinaries()
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	//defer os.RemoveAll(defaultConfig.Base + "/absaoss")
}
