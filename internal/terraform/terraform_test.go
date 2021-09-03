package terraform

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/path"
	"github.com/k0da/tfreg-golang/internal/storage"
	"github.com/k0da/tfreg-golang/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	amd64FileName = "terraform-provider-dummy_1.2.5_linux_amd64.zip"
	arm64FileName = "terraform-provider-dummy_1.2.5_linux_arm64.zip"
)

var platformAmd64 = types.Platform{Os: "linux", Arch: "amd64", FileOrigin: amd64FileName}
var platformArm64 = types.Platform{Os: "linux", Arch: "arm64", FileOrigin: arm64FileName}

var expectedProvider = &Provider{
	path:      getDefaultPath(),
	Platforms: []types.Platform{platformAmd64, platformArm64},
}

var defaultConfig = config.Config{
	Base:        "./../../test_data/target",
	ArtifactDir: "./../../test_data/source",
	Namespace:   "absaoss",
	TargetDir:   "target",
	Branch:      "gh-pages",
	WebRoot:     "/",
}

var greenConfig = config.Config{
	// todo: find a system to store / restore files in test
	Base:        "./../../test_data/target_green",
	ArtifactDir: "./../../test_data/source",
	Namespace:   "absaoss",
	Branch:      "gh-pages",
	WebRoot:     "/",
	Owner:       "absaoss",
	Repository:  "terraform-provider-dummy",
}

func getDefaultPath() *path.Path {
	p, _ := path.NewPath(defaultConfig)
	return p
}

func TestNewProviderParsing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := storage.NewMockStorage(ctrl)
	m.EXPECT().CreatePlatformMetadata(gomock.Any()).AnyTimes().Return(defaultConfig.Base, nil)

	p, _ := path.NewPath(defaultConfig)
	provider, err := NewProvider(p, m)
	require.NoError(t, err)
	assert.Equal(t, expectedProvider.Platforms[0].Arch, provider.Platforms[0].Arch, "expected Architecture %+v, but got: %+v", "amd64", expectedProvider.Platforms[0].Arch)
	assert.Equal(t, expectedProvider.Platforms[1].Arch, provider.Platforms[1].Arch, "expected Architecture %+v, but got: %+v", "arm64", expectedProvider.Platforms[1].Arch)
}

// todo: test corner cases (you create corner cases with wrong data (no files, wrong file names, etc), + wrong config (invalid paths etc...))
func TestVersionFromProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := storage.NewMockStorage(ctrl)
	m.EXPECT().CreatePlatformMetadata(gomock.Any()).AnyTimes().Return(defaultConfig.Base, nil)

	versions := types.Versions{}
	expVersions := types.Versions{}
	p, _ := path.NewPath(defaultConfig)
	provider, err := NewProvider(p, m)
	require.NoError(t, err)
	version := provider.GenerateVersion()
	existing, _ := os.ReadFile(defaultConfig.ArtifactDir + "/existing.json")
	expected, _ := os.ReadFile(defaultConfig.ArtifactDir + "/expected.json")
	err = json.Unmarshal(existing, &versions)
	require.NoError(t, err)
	err = json.Unmarshal(expected, &expVersions)
	require.NoError(t, err)
	versions.Versions = append(versions.Versions, *version)
	assert.Equal(t, expVersions, versions, "VersionsPath doesn't match exp %+v, got: %+v", expVersions, versions)
}

func TestGreenPath(t *testing.T) {
	// init
	path, err := path.NewPath(greenConfig)
	require.NoError(t, err)
	storage, err := storage.NewProvider(path)
	require.NoError(t, err)
	provider, err := NewProvider(path, storage)
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
	// todo: make assertions here
}

func TestMain(m *testing.M) {
	defer os.RemoveAll(greenConfig.Base)
	_ = os.Mkdir(greenConfig.Base+"/absaoss/dummy", 0755)
	_, _ = storage.Copy(greenConfig.ArtifactDir+"/existing.json", greenConfig.Base+"/absaoss/dummy/versions")
	m.Run()
	// todo: prepare folder test_green and clean at the end of the test
}
