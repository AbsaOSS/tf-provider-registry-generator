package config

import (
	"fmt"
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

func TestSaneDefaults(t *testing.T) {
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
	c, err := NewConfig("pages")
	fmt.Printf("AAA %+v", c)
	require.NoError(t, err)
}
