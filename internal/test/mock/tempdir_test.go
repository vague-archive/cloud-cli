package mock_test

import (
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/test/assert"
	"github.com/vaguevoid/cloud-cli/internal/test/mock"
)

//-------------------------------------------------------------------------------------------------

const (
	FirstContent  = "first"
	SecondContent = "second"
	FirstPath     = "path/to/first.txt"
	SecondPath    = "path/to/second.txt"
)

//-------------------------------------------------------------------------------------------------

func TestMockTempDir(t *testing.T) {
	tmp := mock.TempDir(t)

	assert.False(t, tmp.Exists(FirstPath), "preconditions")
	assert.False(t, tmp.Exists(SecondPath), "preconditions")

	tmp.AddTextFile(t, FirstPath, FirstContent)
	tmp.AddTextFile(t, SecondPath, SecondContent)

	assert.True(t, tmp.Exists(FirstPath))
	assert.True(t, tmp.Exists(SecondPath))
}

//-------------------------------------------------------------------------------------------------
