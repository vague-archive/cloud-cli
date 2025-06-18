package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/lib/crypto"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
	"github.com/vaguevoid/cloud-cli/internal/test/mock"
)

//-------------------------------------------------------------------------------------------------

func TestBlake3(t *testing.T) {

	tmp := mock.TempDir(t)
	tmp.AddTextFile(t, "foo.txt", "Hello World")
	f, err := os.Open(filepath.Join(tmp.Dir, "foo.txt"))
	assert.Nil(t, err)
	defer f.Close()

	assert.Equal(t, "41f8394111eb713a22165c46c90ab8f0fd9399c92028fd6d288944b23ff5bf76", crypto.Blake3(f))
	assert.Equal(t, "41f8394111eb713a22165c46c90ab8f0fd9399c92028fd6d288944b23ff5bf76", crypto.Blake3("Hello World"))
	assert.Equal(t, "5342397cf3300914b3a26595fa74115d2020cc209b489d9f4c794d7e64bd5b2b", crypto.Blake3("Yolo"))
	assert.Equal(t, "f890484173e516bfd935ef3d22b912dc9738de38743993cfedf2c9473b3216a4", crypto.Blake3("BLAKE3"))
}

//-------------------------------------------------------------------------------------------------
