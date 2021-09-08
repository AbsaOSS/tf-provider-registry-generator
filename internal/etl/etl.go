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
	location location.ILocation
	storage storage.IStorage
	repo      repo.IRepo
	terraform terraform.ITerraform
}

func NewEtl(c config.Config) (etl *Etl, err error) {
	// todo: Dependency injection
	// todo: constructor accepting particular inputs or ifactory
	etl = new(Etl)
	etl.location, err = location.NewLocation(c)
	if err != nil {
		return
	}
	etl.storage, err = storage.NewStorage(etl.location)
	if err != nil {
		return
	}
	gpg, err := encryption.NewGpg(etl.location)
	if err != nil {
		return
	}
	etl.terraform, err = terraform.NewProvider(etl.location, gpg)
	return
}

func (e *Etl) Run() (err error) {
	err = e.repo.Clone()
	if err != nil {
		return
	}

	downloads, err := e.terraform.GetDownloadInfo()
	if err != nil {
		return
	}

	_, err = e.storage.WritePlatformMetadata(downloads)
	if err != nil {
		return
	}

	versions, err := e.storage.GetVersions()
	if err != nil {
		return
	}

	version := e.terraform.GenerateVersion()
	versions.Versions = append(versions.Versions, *version)
	err = e.storage.WriteVersions(versions)
	if err != nil {
		return
	}
	err = e.storage.SaveBinaries()
	if err != nil {
		return
	}
	err = e.repo.CommitAndPush()
	if err != nil {
		return
	}
	return nil
}
