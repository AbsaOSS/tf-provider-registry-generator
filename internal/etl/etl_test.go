package etl

import (
	"github.com/golang/mock/gomock"
	"github.com/k0da/tfreg-golang/internal/config"
	location2 "github.com/k0da/tfreg-golang/internal/location"
	repo2 "github.com/k0da/tfreg-golang/internal/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

type selector int
const (
	_ selector = 1 << iota
	location
	repo
	terraform
	storage
)

func mocked(t *testing.T,s selector) (b *Batch) {
	b = new(Batch)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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

	case storage:
	}
	return b
}


func TestEtl(t *testing.T) {
		b := mocked(t, repo | storage)
		e := etl.NewEtl(b)
		err := etl.Run()
		assert.NoError(t, err)
}
