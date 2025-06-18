package mock_test

import (
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/test/assert"
	"github.com/vaguevoid/cloud-cli/internal/test/mock"
)

const TestKey = "test-key"
const TestValue = "test-value"

func TestMockKeyring(t *testing.T) {
	keyring := mock.Keyring()
	assert.False(t, keyring.Has(TestKey))

	err := keyring.Set(TestKey, TestValue)
	assert.Nil(t, err)

	value, ok := keyring.Get(TestKey)
	assert.True(t, ok)
	assert.Equal(t, TestValue, value)
	assert.True(t, keyring.Has(TestKey))

	err = keyring.Del(TestKey)
	assert.Nil(t, err)
	assert.False(t, keyring.Has(TestKey))

}
