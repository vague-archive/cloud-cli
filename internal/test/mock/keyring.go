package mock

func Keyring() *MockKeyring {
	return &MockKeyring{
		Store: make(map[string]string),
	}
}

type MockKeyring struct {
	Store map[string]string
}

func (k *MockKeyring) Has(key string) bool {
	_, ok := k.Get(key)
	return ok
}

func (k *MockKeyring) Get(key string) (string, bool) {
	v, ok := k.Store[key]
	return v, ok
}

func (k *MockKeyring) Set(key string, value string) error {
	k.Store[key] = value
	return nil
}

func (k *MockKeyring) Del(key string) error {
	delete(k.Store, key)
	return nil
}
