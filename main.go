package main

import (
	"log"

	"github.com/k0da/tfreg-golang/internal/config"
	location "github.com/k0da/tfreg-golang/internal/location"
	"github.com/k0da/tfreg-golang/internal/repo"
	"github.com/k0da/tfreg-golang/internal/storage"
	"github.com/k0da/tfreg-golang/internal/terraform"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func provider(c config.Config) {
	// init
	location, err := location.NewLocation(c)
	checkError(err)
	storage, err := storage.NewProvider(location)
	checkError(err)
	provider, err := terraform.NewProvider(location)
	checkError(err)
	downloads, err := provider.GetDownloadInfo()
	checkError(err)
	_, err = storage.WritePlatformMetadata(downloads)
	checkError(err)
	versions, err := storage.GetVersions()
	checkError(err)
	version := provider.GenerateVersion()
	versions.Versions = append(versions.Versions, *version)
	err = storage.WriteVersions(versions)
	checkError(err)
	err = storage.SaveBinaries()
	checkError(err)
}

func main() {
	config, err := config.NewConfig("pages")
	checkError(err)
	location, err := location.NewLocation(config)
	checkError(err)
	repo, err := repo.NewGithub(location)
	checkError(err)
	repo.Clone(config)
	provider(config)
	repo.CommitAndPush(config)
}
