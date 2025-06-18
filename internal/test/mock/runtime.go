package mock

func Runtime() *MockRuntime {
	return &MockRuntime{}
}

type MockRuntime struct {
	OpenedURL string
}

func (b *MockRuntime) Open(url string) {
	b.OpenedURL = url
}
