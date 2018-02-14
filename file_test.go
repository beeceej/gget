package gget

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestToFile(t *testing.T) {
	testURLBasedPathExpectSuccess := func() {
		u, _ := url.Parse("http://www.google.com/something/another/this/file.xyz")
		expected := filepath.Join("www.google.com", "something", "another", "this", "file.xyz")
		dest, _ := URLBasedPath(u)

		if expected != dest {
			t.Fail()
		}
	}
	testURLBasedPathExpectSuccess()
}
