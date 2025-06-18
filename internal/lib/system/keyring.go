package system

import "github.com/zalando/go-keyring"

type Keyring interface {
	Has(key string) bool
	Get(key string) (string, bool)
	Set(key string, value string) error
	Del(key string) error
}

func DefaultKeyring(name string) *SystemKeyring {
	return &SystemKeyring{
		Name: name,
	}
}

type SystemKeyring struct {
	Name string
}

func (k *SystemKeyring) Has(key string) bool {
	_, ok := k.Get(key)
	return ok
}

func (k *SystemKeyring) Get(key string) (string, bool) {
	value, err := keyring.Get(k.Name, key)
	if err != nil {
		return "", false
	}
	return value, true
}

func (k *SystemKeyring) Set(key string, value string) error {
	return keyring.Set(k.Name, key, value)
}

func (k *SystemKeyring) Del(key string) error {
	return keyring.Delete(k.Name, key)
}
