package gget

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ToFile returns a strategy for persisting a file
// it will either  write the contents of the file to the destination given,
// or if that destination does not exist, it will return an error
type ToFile struct {
	Base string
	Dst  func(*url.URL) (string, bool)
}

// Handle is a method which saves a result to the file system
func (t *ToFile) Handle(r result) (err error) {
	if dst, ok := t.Dst(r.url); ok {
		return writeFile(r.b, filepath.Join(t.Base, dst))
	}
	return errors.New("Unable to compute destination")
}

func writeFile(b []byte, dst string) (err error) {
	var (
		f *os.File
	)
	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		log.Printf(err.Error())
	}
	if f, err = os.Create(dst); err != nil {
		return fmt.Errorf(fmt.Sprintf("Error creating file %s: %v", dst, err))
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer(b))
	return err
}

// URLBasedPath takes a URL and returns a string representing the path to save to.
//
// for Example the URL https://github.com/beeceej/gget/blob/master/cmd/gget/main.go will return:
// `github.com/beeceej/gget/github.com/beeceej/gget/blob/master/cmd/gget/main.go`
func URLBasedPath(u *url.URL) (path string, ok bool) {
	if u != nil {
		path, ok = filepath.Join(u.Host, filepath.Join(strings.Split(u.Path, string(filepath.Separator))...)), true
	}
	return path, ok
}
