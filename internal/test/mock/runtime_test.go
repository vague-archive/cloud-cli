package mock_test

import (
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/test/assert"
	"github.com/vaguevoid/cloud-cli/internal/test/mock"
)

//-------------------------------------------------------------------------------------------------

const TestURL = "https://play.void.dev"

//-------------------------------------------------------------------------------------------------

func TestMockRuntime(t *testing.T) {
	runtime := mock.Runtime()
	assert.Empty(t, runtime.OpenedURL)
	runtime.Open(TestURL)
	assert.Equal(t, TestURL, runtime.OpenedURL)
}

//-------------------------------------------------------------------------------------------------
