package etl

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k0da/tfreg-golang/internal/config"
	location2 "github.com/k0da/tfreg-golang/internal/location"
	repo2 "github.com/k0da/tfreg-golang/internal/repo"
	storage2 "github.com/k0da/tfreg-golang/internal/storage"
	terraform2 "github.com/k0da/tfreg-golang/internal/terraform"
	"github.com/stretchr/testify/assert"
)

type selector int

const (
	_ selector = 1 << iota
	location
	repo
	terraform
	storage
)

func mock(t *testing.T, c config.Config, s selector) (b *Batch) {
	// get batch based on configuration
	b, _ = NewEtlFactory(c).Get()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// inject mocks into batch
	switch s {
	case location:
		lm := location2.NewMockILocation(ctrl)
		lm.EXPECT().GetConfig().Return(config.Config{}).AnyTimes()
		b.location = lm
	case repo:
		rm := repo2.NewMockIRepo(ctrl)
		rm.EXPECT().Clone().Return(nil).AnyTimes()
		rm.EXPECT().CommitAndPush().Return(nil).AnyTimes()
		b.repo = rm
	case terraform:
		tm := terraform2.NewMockITerraform(ctrl)
		//tm.EXPECT()
		b.terraform = tm
	case storage:
		sm := storage2.NewMockIStorage(ctrl)
		b.storage = sm
	}
	return b
}

func TestEtl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockIFactory(ctrl)
	b := mock(t, config.Config{}, repo|storage)
	m.EXPECT().Get().Return(b).AnyTimes()
	e, err := NewEtl(m)
	assert.NoError(t, err)
	err = e.Run()
	assert.NoError(t, err)
}
