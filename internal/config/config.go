package config

import (
	"fmt"
	"strings"

	"github.com/AbsaOSS/gopkg/env"
)

type Config struct {
	TargetDir      string
	Namespace      string
	ArtifactDir    string
	Branch         string
	WebRoot        string
	Base           string
	Owner          string
	Repository     string
	RepoURL        string
	User           string
	Email          string
	GPGFingerPrint string
	GPGHome        string
}

func NewConfig(base string) (c Config, err error) {
	const targetDir = "TARGET_DIR"
	const artifactsDir = "ARTIFACTS_DIR"
	const namespace = "NAMESPACE"
	const gpgFingerprint = "GPG_FINGERPRINT"
	const branch = "BRANCH"
	const githubRepo = "GITHUB_REPOSITORY"
	const webRoot = "WEB_ROOT"
	const repoURL = "REPO_URL"
	const ghToken = "GITHUB_TOKEN"
	const repo = "REPOSITORY"
	const user = "USERNAME"
	const email = "EMAIL"
	const actor = "GITHUB_ACTOR"
	c = Config{}
	c.TargetDir = env.GetEnvAsStringOrFallback(targetDir, "")
	c.ArtifactDir = env.GetEnvAsStringOrFallback(artifactsDir, "")
	if c.ArtifactDir == "" {
		err = fmt.Errorf("empty %s", c.ArtifactDir)
		return
	}
	c.GPGFingerPrint = env.GetEnvAsStringOrFallback(gpgFingerprint, "")
	if c.GPGFingerPrint == "" {
		err = fmt.Errorf("empty %s", gpgFingerprint)
		return
	}
	c.Branch = env.GetEnvAsStringOrFallback(branch, "gh-pages")
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
	home := env.GetEnvAsStringOrFallback("HOME", "")
	if home == "" {
		err = fmt.Errorf("empty HOME")
	}
	c.GPGHome = home + "/.gnupg"
	ghActor := env.GetEnvAsStringOrFallback(actor, "registry-action")
	c.User = env.GetEnvAsStringOrFallback(user, ghActor)
	c.Email = env.GetEnvAsStringOrFallback(email, ghActor + "@users.noreply.github.com")

	return
}
