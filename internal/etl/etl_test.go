package etl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AbsaOSS/tf-provider-registry-generator/internal/config"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/repo"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/storage"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var greenConfig = config.Config{
	// todo: find a system to store / restore files in test
	Base:        "./../../test_data/target_etl",
	ArtifactDir: "./../../test_data/source",
	Namespace:   "absaoss",
	Branch:      "gh-pages",
	WebRoot:     "/",
	Owner:       "absaoss",
	Repository:  "terraform-provider-dummy",
}

func TestEtl(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rm := repo.NewMockIRepo(ctrl)
	asset := &types.FileAsset{
		SHASum: "shasum",
		SHASig: "shasig",
		Download: map[string]string{
			"linux_amd64": "linux_x86",
			"linux_arm64": "linux_arm",
		},
	}
	rm.EXPECT().Clone().Return(nil).AnyTimes()
	rm.EXPECT().CommitAndPush().Return(nil).AnyTimes()
	rm.EXPECT().GetAssets("1.2.3").Return(asset).AnyTimes()
	b, _ := NewEtlFactory(greenConfig).Get()
	b.repo = rm
	f := NewMockIFactory(ctrl)
	f.EXPECT().Get().Return(b, nil).AnyTimes()
	e, err := NewEtl(f)
	require.NoError(t, err)

	// act
	err = e.Run()
	assert.NoError(t, err)
	assert.True(t, exists(greenConfig, "/absaoss/dummy/1.2.5/download/linux/amd64"))
	assert.True(t, exists(greenConfig, "/absaoss/dummy/1.2.5/download/linux/arm64"))
	assert.True(t, exists(greenConfig, "/terraform.json"))
}

func exists(config config.Config, subpath string) bool {
	if _, err := os.Stat(config.Base + subpath); os.IsNotExist(err) {
		return false
	}
	return true
}

func TestMain(m *testing.M) {
	defer os.RemoveAll(greenConfig.Base)
	err := os.Mkdir(greenConfig.Base, 0755)
	if err != nil {
		os.Exit(1)
	}
	_, _ = storage.Copy(greenConfig.ArtifactDir+"/existing.json", greenConfig.Base+"/absaoss/dummy/versions")
	m.Run()
	// todo: prepare folder test_green and clean at the end of the test
}
