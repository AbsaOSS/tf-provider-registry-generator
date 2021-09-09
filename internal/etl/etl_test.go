package etl

import (
	"github.com/k0da/tfreg-golang/internal/storage"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/repo"
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
	rm.EXPECT().Clone().Return(nil).AnyTimes()
	rm.EXPECT().CommitAndPush().Return(nil).AnyTimes()
	b,_ := NewEtlFactory(greenConfig).Get()
	b.repo = rm
	f := NewMockIFactory(ctrl)
	f.EXPECT().Get().Return(b,nil).AnyTimes()
	e, err := NewEtl(f)
	require.NoError(t, err)

	// act
	err = e.Run()
	assert.NoError(t, err)
	assert.True(t, exists(greenConfig, "/absaoss/dummy/1.2.5/download/linux/amd64"))
	assert.True(t, exists(greenConfig, "/absaoss/dummy/1.2.5/download/linux/arm64"))
	assert.True(t, exists(greenConfig, "/terraform.json"))
}

func exists(config config.Config, subpath string) bool{
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