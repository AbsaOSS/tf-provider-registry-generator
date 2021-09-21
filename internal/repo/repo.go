package repo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/AbsaOSS/gopkg/shell"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/location"
	"github.com/AbsaOSS/tf-provider-registry-generator/internal/types"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const commitMsg = "Generate Terraform registry"

type IRepo interface {
	Clone() (err error)
	CommitAndPush() (err error)
	GetAssets(string, []*github.RepositoryRelease) (asset *types.FileAsset, err error)
	GetReleases() ([]*github.RepositoryRelease, error)
}

type Github struct {
	location location.ILocation
}

func NewGithub(location location.ILocation) (g *Github, err error) {
	if location == nil {
		return nil, fmt.Errorf("nil location")
	}
	g = &Github{
		location: location,
	}
	return
}

func filterRelease(releases []*github.RepositoryRelease, v string) (assets []github.ReleaseAsset, err error) {
	for _, r := range releases {
		if strings.Contains(*r.TagName, v) {
			assets = r.Assets
		}
	}
	if len(assets) == 0 {
		err = fmt.Errorf("No assets found for release %s", v)
	}

	return assets, err
}
func isArtifact(artifact, asset string) bool {
	return artifact == asset
}
func isNotSet(s string) bool {
	return s == ""
}
func isSHASum(s string) bool {
	return strings.HasSuffix(s, "SHA256SUMS")
}
func isSHASig(s string) bool {
	return strings.HasSuffix(s, "SHA256SUMS.sig")
}
func (g *Github) GetReleases() ([]*github.RepositoryRelease, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.location.GetConfig().Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)
	releases, _, err := ghClient.Repositories.ListReleases(ctx, g.location.GetConfig().Owner, g.location.GetConfig().Repository, nil)
	return releases, err
}

func (g *Github) GetAssets(v string, releases []*github.RepositoryRelease) (files *types.FileAsset, err error) {
	artifacts := g.location.GetArtifacts()
	files = new(types.FileAsset)
	files.Download = make(map[string]string)
	assets, err := filterRelease(releases, v)
	if err != nil {
		return nil, err
	}
	for _, asset := range assets {
		for _, artifact := range artifacts {
			osArch := fmt.Sprintf("%s_%s", artifact.Os, artifact.Arch)
			if isArtifact(*asset.Name, artifact.File) {
				if isNotSet(files.Download[osArch]) {
					files.Download[osArch] = *asset.URL
				}
			}
			if isSHASum(*asset.Name) {
				files.SHASum = *asset.URL
			}
			if isSHASig(*asset.Name) {
				files.SHASig = *asset.URL
			}
		}
	}
	return files, err
}

func (g *Github) Clone() (err error) {
	// todo: use location instead of config
	args := []string{"clone", "--branch", g.location.GetConfig().Branch, g.location.GetConfig().RepoURL, g.location.GetConfig().Base}
	err = runGitCmd(args, "")
	return
}

func (g *Github) CommitAndPush() (err error) {
	gitUser := []string{"config", "user.name", g.location.GetConfig().User}
	err = runGitCmd(gitUser, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitEmail := []string{"config", "user.email", g.location.GetConfig().Email}
	err = runGitCmd(gitEmail, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitSetRemote := []string{"remote", "set-url", "origin", g.location.GetConfig().RepoURL}
	err = runGitCmd(gitSetRemote, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitAdd := []string{"add", "./"}
	err = runGitCmd(gitAdd, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitCommit := []string{"commit", "-m", commitMsg}
	err = runGitCmd(gitCommit, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	gitPush := []string{"push", "origin", g.location.GetConfig().Branch}
	err = runGitCmd(gitPush, g.location.GetConfig().Base)
	if err != nil {
		return
	}

	return
}

func runGitCmd(args []string, dir string) (err error) {
	command := shell.Command{
		Command:    "git",
		Args:       args,
		WorkingDir: dir,
	}
	cmdLog, err := shell.Execute(command)
	if err != nil {
		return
	}
	log.Printf("run git %+v: %s", args, cmdLog)
	return
}
