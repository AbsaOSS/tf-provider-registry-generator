package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/k0da/tfreg-golang/internal/terraform"
	"github.com/stretchr/testify/assert"
)


func TestGenerateDownloadPath(t *testing.T) {
	path, err := os.MkdirTemp("./", "test-pages-")
	assert.NoError(t, err)
	err = os.Setenv("NAMESPACE", "absaoss")
	assert.NoError(t, err)
	prepareDownloadDir(expectedProvider)
	t.Logf("got path %s", path)
}
