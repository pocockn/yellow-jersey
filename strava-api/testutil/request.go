package testutil

import "io"

// NoopCloser implements the io.ReadCloser interface
// Is used in HTTP request tests
type NoopCloser struct {
	io.Reader
}

// Close implements the io.ReadCloser interface to allow us to use NoopCloser as a HTTP body.
func (NoopCloser) Close() error { return nil }
