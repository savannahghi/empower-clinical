package mock

import (
	"context"
	"net/http"
)

// FakeISCExtension is an mock of the ISCExtension
type FakeISCExtension struct {
	MockMakeRequestFn func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
}

// NewFakeISCExtensionMock initializes a new instance of extension mock
func NewFakeISCExtensionMock() *FakeISCExtension {
	return &FakeISCExtension{
		MockMakeRequestFn: func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
			return nil, nil
		},
	}
}

// MakeRequest mocks the implementation of MakeRequest method
func (b *FakeISCExtension) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	return b.MockMakeRequestFn(ctx, method, path, body)
}
