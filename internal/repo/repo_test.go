package repo

import (
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/config"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/location"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/types"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSHASig(t *testing.T) {
	file := "terraform-provider-hashicups_0.3.2_SHA256SUMS.sig"
	assert.True(t, isSHASig(file), "Should be true")
}
func TestSHASums(t *testing.T) {
	file := "terraform-provider-hashicups_0.3.2_SHA256SUMS"
	assert.True(t, isSHASum(file), "Should be true")
}

func TestGetAssets(t *testing.T) {
	c := config.Config{
		ArtifactDir: "../../test_data/source",
	}
	l, err := location.NewLocation(c)
	require.NoError(t, err)
	g := &Github{
		location: l,
	}
	assets := []github.ReleaseAsset{
		{
			Name: github.String("terraform-provider-dummy_1.2.5_linux_amd64.zip"),
			URL:  github.String("http://api.github/release/asset/1"),
		},
		{
			Name: github.String("terraform-provider-dummy_1.2.5_linux_arm64.zip"),
			URL:  github.String("http://api.github/release/asset/2"),
		},
		{
			Name: github.String("terraform-provider-dummy_1.2.5_SHA256SUMS.sig"),
			URL:  github.String("http://api.github/release/asset/3"),
		},
		{
			Name: github.String("terraform-provider-dummy_1.2.5_SHA256SUMS"),
			URL:  github.String("http://api.github/release/asset/4"),
		},
	}
	releases := []*github.RepositoryRelease{{
		TagName: github.String("v1.2.5"),
		Assets:  assets,
	}}
	expFileAsset := &types.FileAsset{
		SHASum: "http://api.github/release/asset/4",
		SHASig: "http://api.github/release/asset/3",
		Download: map[string]string{
			"linux_amd64": "http://api.github/release/asset/1",
			"linux_arm64": "http://api.github/release/asset/2",
		},
	}
	got, err := g.GetAssets("1.2.5", releases)
	require.NoError(t, err)
	assert.Equal(t, expFileAsset, got, "expected %s, but got: %s", expFileAsset, got)
}
