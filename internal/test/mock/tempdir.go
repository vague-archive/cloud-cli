package mock

import (
	"os"
	"path/filepath"
	"testing"
)

func TempDir(t *testing.T) *MockTempDir {
	return &MockTempDir{
		Dir: t.TempDir(),
	}
}

type MockTempDir struct {
	Dir string
}

func (tmp *MockTempDir) AddTextFile(t *testing.T, name string, content string) {
	tmp.AddFile(t, name, []byte(content))
}

func (tmp *MockTempDir) AddFile(t *testing.T, name string, content []byte) {
	fullPath := filepath.Join(tmp.Dir, name)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func (tmp *MockTempDir) Exists(name string) bool {
	fullPath := filepath.Join(tmp.Dir, name)
	_, err := os.Stat(fullPath)
	return err == nil
}
