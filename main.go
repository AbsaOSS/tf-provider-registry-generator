package main

import (
	"log"

	"github.com/k0da/tfreg-golang/internal/config"
	"github.com/k0da/tfreg-golang/internal/github"
	pather "github.com/k0da/tfreg-golang/internal/location"
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
	location, err := pather.NewLocation(c)
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
	github.Clone(config)
	provider(config)
	github.CommitAndPush(config)
}
