package etl

import (
	"fmt"
	"log"

	"github.com/k0da/tfreg-golang/internal/location"
	"github.com/k0da/tfreg-golang/internal/repo"
	"github.com/k0da/tfreg-golang/internal/storage"
	"github.com/k0da/tfreg-golang/internal/terraform"
)

type IEtl interface {
	Run() error
}

type Etl struct {
	location  location.ILocation
	storage   storage.IStorage
	repo      repo.IRepo
	terraform terraform.ITerraform
}

// NewEtl accepts factory which retreives valid Batch or error. The reason it accept factory instead of Batch is
// 1. extensibility (open closed principle, SOLID)
// 2. I don't need to make extra validations of particular Batch fields
// 3. easier testing / mocking
func NewEtl(f IFactory) (etl *Etl, err error) {
	etl = new(Etl)
	if f == nil {
		err = fmt.Errorf("nil IFactory")
		return
	}
	b, err := f.Get()
	if err != nil {
		return
	}
	etl.storage = b.storage
	etl.location = b.location
	etl.terraform = b.terraform
	etl.repo = b.repo
	return
}

func (e *Etl) Run() (err error) {
	log.Printf("Doing clone")
	err = e.repo.Clone()
	if err != nil {
		return
	}

	log.Printf("Generate terraform.json")
	err = e.terraform.GenerateTerraformJSON()
	if err != nil {
		return
	}

	log.Printf("Generate download info")
	downloads, err := e.terraform.GetDownloadInfo()
	if err != nil {
		return
	}

	log.Printf("Writing download metadata")
	_, err = e.storage.WritePlatformMetadata(downloads)
	if err != nil {
		return
	}

	log.Printf("Fetching versions")
	versions, err := e.storage.GetVersions()
	if err != nil {
		return
	}

	log.Printf("Generate version")
	version := e.terraform.GenerateVersion()
	versions.Versions = append(versions.Versions, *version)
	log.Printf("Merging version to versions and writing it")
	err = e.storage.WriteVersions(versions)
	if err != nil {
		return
	}

	log.Printf("Pushing")
	err = e.repo.CommitAndPush()


	return nil
}
