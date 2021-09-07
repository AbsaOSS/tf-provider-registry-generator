package etl

import (
	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/encryption"
	"github.com/k0da/tfreg-golang/internal/location"
	"github.com/k0da/tfreg-golang/internal/repo"
	"github.com/k0da/tfreg-golang/internal/storage"
	"github.com/k0da/tfreg-golang/internal/terraform"
)

type IEtl interface {
	Run() error
}

type Etl struct {
	l location.ILocation
	s storage.IStorage
	r repo.IRepo
	t terraform.ITerraform
}

func NewEtl(c config.Config) (etl *Etl, err error) {
	// todo: Dependency injection
	etl = new(Etl)
	etl.l, err = location.NewLocation(c)
	if err != nil {
		return
	}
	etl.s, err = storage.NewStorage(etl.l)
	if err != nil {
		return
	}
	gpg, err := encryption.NewGpg(etl.l)
	if err != nil {
		return
	}
	etl.t, err = terraform.NewProvider(etl.l, gpg)
	return
}

func (e *Etl) Run() error {
	// todo: implement workflow here
	return nil
}
