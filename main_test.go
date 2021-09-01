package main

import (
	"testing"
	"os"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

const fileName = "terraform-provider-dummy_1.2.5_linux_amd64.zip"
var platform = Platform{ Os: "linux", Arch: "amd64"}
var expectedProvider = &Provider{
		Name: "dummy",
		Version: "1.2.5",
		Files: []string{fileName},
		Platforms: []Platform{platform},
	}

func TestProviderParsing(t *testing.T) {
	provider := parseProvider(fileName)
	provider.updatePlatform(getOs(fileName), getArch(fileName))
	gotName := provider.Name
	gotVer := provider.Version
	assert.Equal(t, expectedProvider.Name, gotName, "expected %s, but got: %s", expectedProvider.Name, gotName)
	assert.Equal(t, expectedProvider.Version, gotVer, "expected %s, but got: %s", expectedProvider.Version, gotVer)
	assert.Equal(t, expectedProvider, provider, "expected Provider %+v, but got: %+v", expectedProvider, provider)
}

func TestVersionFromProvider(t *testing.T){
	versions := Versions{}
	expVersions := Versions{}
	version := expectedProvider.generateVersion()
	existing, _ := os.ReadFile("./fixtures/existing.json")
	expected, _ := os.ReadFile("./fixtures/expected.json")
	err := json.Unmarshal(existing, &versions)
	assert.NoError(t, err)
	err = json.Unmarshal(expected, &expVersions)
	assert.NoError(t, err)
	versions.Versions = append(versions.Versions, *version)
	assert.Equal(t, expVersions, versions, "Versions doesn't match exp %+v, got: %+v", expVersions, versions)
}

func TestGenerateDownloadPath(t *testing.T){
	path, err := os.MkdirTemp("./", "test-pages-")
	assert.NoError(t, err)
	err = os.Setenv("NAMESPACE", "absaoss")
	assert.NoError(t, err)
	prepareDownloadDir(expectedProvider)
	t.Logf("got path %s", path)
}
