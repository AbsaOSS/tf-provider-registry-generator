package config

import (
	"fmt"
	"strings"

	"github.com/AbsaOSS/gopkg/env"
)

type Config struct {
	TargetDir   string
	Namespace   string
	ArtifactDir string
	Branch      string
	WebRoot     string
	Base        string
	Owner       string
	Repository  string
	RepoURL     string
	User        string
	Email       string
	GPGKeyID    string
	GPGArmor    string
}

func NewConfig(base string) (c Config, err error) {
	const targetDir = "TARGET_DIR"
	const artifactsDir = "ARTIFACTS_DIR"
	const namespace = "NAMESPACE"
	const gpgArmor = "GPG_ASCII_ARMOR"
	const gpgKeyID = "GPG_KEYID"
	const branch = "BRANCH"
	const webRoot = "WEB_ROOT"
	const ghToken = "TOKEN"
	const repo = "REPOSITORY"
	const user = "USERNAME"
	const email = "EMAIL"
	const actor = "GITHUB_ACTOR"
	const githubRepo = "GITHUB_REPOSITORY"
	c = Config{}
	// with defaults
	c.Branch = env.GetEnvAsStringOrFallback(branch, "gh-pages")
	c.ArtifactDir = env.GetEnvAsStringOrFallback(artifactsDir, "dist")

	// mandatory
	c.TargetDir = env.GetEnvAsStringOrFallback(targetDir, "")
	c.GPGKeyID = env.GetEnvAsStringOrFallback(gpgKeyID, "")
	if c.GPGKeyID == "" {
		err = fmt.Errorf("empty %s", gpgKeyID)
		return
	}
	c.GPGArmor = env.GetEnvAsStringOrFallback(gpgArmor, "")
	if c.GPGArmor == "" {
		err = fmt.Errorf("empty %s", gpgArmor)
		return
	}
	ghRepo := env.GetEnvAsStringOrFallback(githubRepo, "")

	orgRepo := strings.Split(ghRepo, "/")
	if len(orgRepo) != 2 {
		err = fmt.Errorf("failed to parse %s", ghRepo)
		return
	}
	c.Owner = orgRepo[0]
	c.Repository = orgRepo[1]

	targetRepo := env.GetEnvAsStringOrFallback(repo, "")
	if targetRepo == "" {
		err = fmt.Errorf("empty %s", repo)
		return
	}

	c.WebRoot = env.GetEnvAsStringOrFallback(webRoot, "/")
	c.Namespace = env.GetEnvAsStringOrFallback(namespace, c.Owner)
	token := env.GetEnvAsStringOrFallback(ghToken, "")
	if token == "" {
		err = fmt.Errorf("empty token")
		return
	}
	c.RepoURL = fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s", token, c.Owner, targetRepo)
	if base == "" {
		err = fmt.Errorf("empty base")
		return
	}
	c.Base = base
	ghActor := env.GetEnvAsStringOrFallback(actor, "registry-action")
	c.User = env.GetEnvAsStringOrFallback(user, ghActor)
	c.Email = env.GetEnvAsStringOrFallback(email, ghActor+"@users.noreply.github.com")

	return
}
