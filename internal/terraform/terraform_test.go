package terraform

import (
	"encoding/json"
	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	amd64FileName = "terraform-provider-dummy_1.2.5_linux_amd64.zip"
	arm64FileName = "terraform-provider-dummy_1.2.5_linux_arm64.zip"
)

var platformAmd64 = Platform{Os: "linux", Arch: "amd64", fileName: amd64FileName}
var platformArm64 = Platform{Os: "linux", Arch: "arm64", fileName: arm64FileName}

var expectedProvider = &TerraformProvider{
	Name:      "dummy",
	Version:   "1.2.5",
	Platforms: []Platform{platformAmd64, platformArm64},
}

var defaultConfig = config.Config{
	Base:  "./test_pages",
	ArtifactDir: "./../../test_data/source",
	Namespace: "absaoss",
	TargetDir: "target",
	Branch: "gh-pages",
	WebRoot: "/",
}

func TestNewProviderParsing(t *testing.T) {
	provider,err := NewProvider(defaultConfig)
	require.NoError(t, err)
	assert.Equal(t, expectedProvider.Name, provider.Name, "expected %s, but got: %s", expectedProvider.Name, provider.Name)
	assert.Equal(t, expectedProvider.Version, provider.Version, "expected %s, but got: %s", expectedProvider.Version, provider.Version)
	assert.Equal(t, expectedProvider, provider, "expected Provider %+v, but got: %+v", expectedProvider, provider)
	assert.Equal(t, expectedProvider.Platforms[0].Arch, provider.Platforms[0].Arch, "expected Architecture %+v, but got: %+v", "amd64", expectedProvider.Platforms[0].Arch)
	assert.Equal(t, expectedProvider.Platforms[1].Arch, provider.Platforms[1].Arch, "expected Architecture %+v, but got: %+v", "arm64", expectedProvider.Platforms[1].Arch)
}
// todo: test corner cases (you create corner cases with wrong data (no files, wrong file names, etc), + wrong config (invalid paths etc...))


func TestVersionFromProvider(t *testing.T) {
	versions := Versions{}
	expVersions := Versions{}
	provider,err := NewProvider(defaultConfig)
	require.NoError(t, err)
	version := provider.GenerateVersion()
	existing, _ := os.ReadFile(defaultConfig.ArtifactDir+"/existing.json")
	expected, _ := os.ReadFile(defaultConfig.ArtifactDir+"/expected.json")
	err = json.Unmarshal(existing, &versions)
	require.NoError(t, err)
	err = json.Unmarshal(expected, &expVersions)
	require.NoError(t, err)
	versions.Versions = append(versions.Versions, *version)
	assert.Equal(t, expVersions, versions, "Versions doesn't match exp %+v, got: %+v", expVersions, versions)
}
