package terraform

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/AbsaOSS/tf-provider-registry-generator/internal/config"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/location"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/storage"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	amd64FileName = "terraform-provider-dummy_1.2.5_linux_amd64.zip"
	arm64FileName = "terraform-provider-dummy_1.2.5_linux_arm64.zip"
)

var platformAmd64 = types.Platform{Os: "linux", Arch: "amd64", FileOrigin: amd64FileName}
var platformArm64 = types.Platform{Os: "linux", Arch: "arm64", FileOrigin: arm64FileName}

var expectedProvider = &TerraformProvider{
	location:  getDefaultPath(),
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

func prepare() {
	_ = os.Mkdir(greenConfig.Base+"/absaoss/dummy", 0755)
	_, _ = storage.Copy(greenConfig.ArtifactDir+"/existing.json", greenConfig.Base+"/absaoss/dummy/versions")
}

func cleanup() {
	os.RemoveAll(greenConfig.Base)
}

func getDefaultPath() *location.Location {
	p, _ := location.NewLocation(defaultConfig)
	return p
}

func TestNewProviderParsing(t *testing.T) {
	defer cleanup()
	prepare()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := storage.NewMockIStorage(ctrl)
	m.EXPECT().WritePlatformMetadata(gomock.Any()).AnyTimes().Return(defaultConfig.Base, nil)

	l, _ := location.NewLocation(defaultConfig)
	provider, err := NewTerraformProvider(l)
	require.NoError(t, err)
	assert.Equal(t, expectedProvider.Platforms[0].Arch, provider.Platforms[0].Arch, "expected Architecture %+v, but got: %+v", "amd64", expectedProvider.Platforms[0].Arch)
	assert.Equal(t, expectedProvider.Platforms[1].Arch, provider.Platforms[1].Arch, "expected Architecture %+v, but got: %+v", "arm64", expectedProvider.Platforms[1].Arch)
}

// todo: test corner cases (you create corner cases with wrong data (no files, wrong file names, etc), + wrong config (invalid paths etc...))
func TestVersionFromProvider(t *testing.T) {
	defer cleanup()
	prepare()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := storage.NewMockIStorage(ctrl)
	m.EXPECT().WritePlatformMetadata(gomock.Any()).AnyTimes().Return(defaultConfig.Base, nil)

	versions := types.Versions{}
	expVersions := types.Versions{}
	l, _ := location.NewLocation(defaultConfig)
	provider, err := NewTerraformProvider(l)
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
	defer cleanup()
	prepare()
	location, err := location.NewLocation(greenConfig)
	require.NoError(t, err)
	storage, err := storage.NewStorage(location)
	require.NoError(t, err)
	provider, err := NewTerraformProvider(location)
	require.NoError(t, err)
	assets := &types.FileAsset{
		SHASum: "http://foo/shasum",
		SHASig: "http://foo/shasig",
		Download: map[string]string{
			"linux_amd64": "http://linux_amd",
			"linux_arm64": "http://darwin",
		},
	}
	downloads, err := provider.GetDownloadInfo(assets)
	require.NoError(t, err)
	_, err = storage.WritePlatformMetadata(downloads)
	require.NoError(t, err)
	versions, err := storage.GetVersions()
	require.NoError(t, err)
	version := provider.GenerateVersion()
	versions.Versions = append(versions.Versions, *version)
	err = storage.WriteVersions(versions)
	require.NoError(t, err)
	// todo: make assertions here
}
