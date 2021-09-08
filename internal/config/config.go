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
	const targetDir = "INPUT_TARGET_DIR"
	const artifactsDir = "INPUT_ARTIFACTS_DIR"
	const namespace = "INPUT_NAMESPACE"
	const gpgArmor = "INPUT_GPG_ASCII_ARMOR"
	const gpgKeyID = "INPUT_GPG_KEYID"
	const branch = "INPUT_BRANCH"
	const webRoot = "WEB_ROOT"
	const repoURL = "REPO_URL"
	const ghToken = "INPUT_TOKEN"
	const repo = "INPUT_REPOSITORY"
	const user = "INPUT_USERNAME"
	const email = "INPUT_EMAIL"
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
	c.Repository = env.GetEnvAsStringOrFallback(repo, orgRepo[1])

	c.WebRoot = env.GetEnvAsStringOrFallback(webRoot, "/")
	c.Namespace = env.GetEnvAsStringOrFallback(namespace, c.Owner)
	token := env.GetEnvAsStringOrFallback(ghToken, "")
	if token == "" {
		err = fmt.Errorf("empty token")
		return
	}
	fallbackRepo := fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s", token, c.Owner, c.Repository)
	c.RepoURL = env.GetEnvAsStringOrFallback(repoURL, fallbackRepo)
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
