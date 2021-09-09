package etl

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/repo"
	"github.com/stretchr/testify/assert"
)

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
}
