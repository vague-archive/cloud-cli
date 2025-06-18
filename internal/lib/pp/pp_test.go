package pp_test

import (
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/domain/account"
	"github.com/vaguevoid/cloud-cli/internal/lib/pp"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
)

func TestJson(t *testing.T) {
	user := &account.User{
		ID:   42,
		Name: "Jake",
	}
	str := pp.JSON(user)
	assert.Equal(t, pp.Dedent(`
    {
      "id": 42,
      "name": "Jake"
    }`), str)
}

//-------------------------------------------------------------------------------------------------

func TestDedent(t *testing.T) {
	assert.Equal(t, "", pp.Dedent(""))
	assert.Equal(t, "", pp.Dedent("     "))
	assert.Equal(t, "hello world", pp.Dedent("hello world"))
	assert.Equal(t, "hello world", pp.Dedent("  hello world"))
	assert.Equal(t, "hello world", pp.Dedent("    hello world"))
	assert.Equal(t, "hello world", pp.Dedent("\thello world"))
	assert.Equal(t, "{\n  id: 42,\n  name: \"jake\"\n}", pp.Dedent(`
  {
    id: 42,
    name: "jake"
  }`))

	assert.Equal(t, "  {\n  id: 42,\n  name: \"jake\"\n}", pp.Dedent(`
  {
  id: 42,
  name: "jake"
}`))

	assert.Equal(t, "{\n  id: 42,\n  name: \"jake\"\n  }", pp.Dedent(`
{
  id: 42,
  name: "jake"
  }`))

	assert.Equal(t, "{\n    id: 42,\n    name: \"jake\"\n  }", pp.Dedent(`{
    id: 42,
    name: "jake"
  }`))
}

//-------------------------------------------------------------------------------------------------
