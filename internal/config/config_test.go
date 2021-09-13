package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanup() {
	os.Unsetenv(ghToken)
	os.Unsetenv(githubRepo)
	os.Unsetenv(actor)
	os.Unsetenv(gpgArmor)
	os.Unsetenv(gpgKeyID)
}

func setupBasic() {
	os.Setenv(ghToken, "123")
	os.Setenv(gpgArmor, "testing-key\n\n")
	os.Setenv(gpgKeyID, "12345")
	os.Setenv(githubRepo, "foo/bar")
}

func TestSaneDefaults(t *testing.T) {
	defer cleanup()
	os.Setenv(githubRepo, "foo/bar")
	os.Setenv(actor, "testing")
	os.Setenv(ghToken, "123")
	os.Setenv(gpgArmor, "testing-key\n\n")
	os.Setenv(gpgKeyID, "12345")
	expConfig := Config{
		TargetDir:   "",
		Namespace:   "foo",
		ArtifactDir: "dist",
		Branch:      "gh-pages",
		WebRoot:     "/",
		Base:        "pages",
		Owner:       "foo",
		Repository:  "bar",
		RepoURL:     "https://x-access-token:123@github.com/foo/bar",
		User:        "testing",
		Email:       "testing@users.noreply.github.com",
		GPGKeyID:    "12345",
		GPGArmor:    "testing-key\n\n",
	}
	c, err := NewConfig("pages")
	require.NoError(t, err)
	assert.Equal(t, c, expConfig, "Expecting config %+v, but got: %+v", expConfig, c)
}

func TestRequireToken(t *testing.T) {
	defer cleanup()
	_, err := NewConfig("pages")
	require.Error(t, err)
}
func TestRequireGPGKey(t *testing.T) {
	defer cleanup()
	os.Setenv(ghToken, "123")
	_, err := NewConfig("pages")
	require.Error(t, err)
}
func TestRequireGPGArmor(t *testing.T) {
	defer cleanup()
	os.Setenv(ghToken, "123")
	os.Setenv(gpgKeyID, "keyid")
	_, err := NewConfig("pages")
	require.Error(t, err)
}

func TestSetupOverrides(t *testing.T) {
	setupBasic()
	expectedRepo := "repoOverride"
	expectedRepoURL := "https://x-access-token:123@github.com/foo/" + expectedRepo
	os.Setenv(repo, expectedRepo)
	c, err := NewConfig("pages")
	require.NoError(t, err)
	assert.Equal(t, expectedRepoURL, c.RepoURL, "expected %s, got: %s", expectedRepoURL, c.RepoURL)
}
