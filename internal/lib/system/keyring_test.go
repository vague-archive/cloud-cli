package system_test

import (
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/lib/system"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
)

const TestKey = "test-key"
const TestValue = "test-value"
const TestDomain = "https://test.void.dev/"

//-------------------------------------------------------------------------------------------------

func TestDefaultKeyringName(t *testing.T) {
	keyring := system.DefaultKeyring(TestDomain)
	assert.NotNil(t, keyring)
	assert.Equal(t, TestDomain, keyring.Name)
}

//-------------------------------------------------------------------------------------------------

func TestKeyring(t *testing.T) {
	keyring := system.DefaultKeyring(TestDomain)

	assert.NotNil(t, keyring)
	assert.Equal(t, TestDomain, keyring.Name)

	err := keyring.Set("canary", "coalmine")
	if err != nil {
		t.Skipf("Skipping test: keyring not available: %v", err)
	}

	assert.False(t, keyring.Has(TestKey))
	err = keyring.Set(TestKey, TestValue)
	assert.Nil(t, err)

	value, ok := keyring.Get(TestKey)
	assert.True(t, ok)
	assert.Equal(t, TestValue, value)
	assert.True(t, keyring.Has(TestKey))

	err = keyring.Del(TestKey)
	assert.Nil(t, err)
	assert.False(t, keyring.Has(TestKey))
}

//-------------------------------------------------------------------------------------------------
