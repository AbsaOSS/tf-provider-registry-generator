package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/AbsaOSS/tf-provider-registry-generator/internal/location"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/types"
)

// Provider retrieving proper location
type Storage struct {
	location location.ILocation
}

type IStorage interface {
	WritePlatformMetadata(download []types.Download) (path string, err error)
	GetVersions() (v types.Versions, err error)
	WriteVersions(v types.Versions) (err error)
}

func NewStorage(l location.ILocation) (provider *Storage, err error) {
	provider = new(Storage)
	provider.location = l
	return
}

func (s *Storage) WritePlatformMetadata(downloads []types.Download) (path string, err error) {
	var b []byte
	for _, d := range downloads {
		dir := s.location.DownloadsPath() + "/" + d.Os
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return
		}
		path = dir + "/" + d.Arch
		b, err = json.Marshal(d)
		if err != nil {
			return
		}
		err = os.WriteFile(path, b, 0644)
		if err != nil {
			return
		}
	}
	return
}

// GetVersions takes versions.json and retreives Versions struct
// if file doesn't exists, return empty Versions slice
func (s *Storage) GetVersions() (v types.Versions, err error) {
	v = types.Versions{}
	if _, err := os.Stat(s.location.VersionsPath()); os.IsNotExist(err) {
		return v, nil
	}
	data, err := os.ReadFile(s.location.VersionsPath())
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

// WriteVersions stores versions.json
func (s *Storage) WriteVersions(v types.Versions) (err error) {
	if len(v.Versions) == 0 {
		err = fmt.Errorf("empty versions")
		return
	}
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	err = os.WriteFile(s.location.VersionsPath(), data, 0644)
	return
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
