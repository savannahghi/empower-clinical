package hapifhir

import (
	"bytes"
	"io"
	"net/http"
)

// jsonPatchContentType is what a FHIR server expects when the PATCH body is a
// JSON Patch document rather than a resource.
const jsonPatchContentType = "application/json-patch+json"

// patchContentTypeTransport corrects the Content-Type on JSON Patch requests.
//
// hapi-fhir-go's FHIRPathPatch builds a JSON Patch array, but its setHeaders
// applies Content-Type: application/fhir+json to every outbound request. HAPI
// then tries to parse the array as a FHIR resource and rejects it with
// "Content does not appear to be FHIR JSON, first non-whitespace character
// was: '['". That affects every patching call in this service - ending an
// encounter or episode, and patching allergies, compositions, tasks,
// conditions, medication requests and patients.
//
// Correcting it here rather than at each call site means the fix holds for all
// of them, and disappears cleanly if the library is fixed upstream: a body that
// is not a JSON array is passed through untouched.
type patchContentTypeTransport struct {
	base http.RoundTripper
}

// NewPatchContentTypeTransport wraps a round tripper, defaulting to
// http.DefaultTransport when none is supplied.
func NewPatchContentTypeTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}

	return &patchContentTypeTransport{base: base}
}

func (t *patchContentTypeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Method != http.MethodPatch || r.Body == nil {
		return t.base.RoundTrip(r)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	_ = r.Body.Close()

	// Restore the body regardless of what we decide about the header.
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))

	if isJSONArray(body) {
		// Clone rather than mutate: the caller's request may be reused.
		clone := r.Clone(r.Context())
		clone.Body = io.NopCloser(bytes.NewReader(body))
		clone.ContentLength = int64(len(body))
		clone.Header.Set("Content-Type", jsonPatchContentType)

		return t.base.RoundTrip(clone)
	}

	return t.base.RoundTrip(r)
}

// isJSONArray reports whether the first non-whitespace byte opens an array,
// which is how a JSON Patch document is distinguished from a FHIR resource.
func isJSONArray(body []byte) bool {
	for _, b := range body {
		switch b {
		case ' ', '\t', '\r', '\n':
			continue
		case '[':
			return true
		default:
			return false
		}
	}

	return false
}
