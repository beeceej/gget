package gget

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jofo8948/gget/src/strategy"
)

// GGet provides a way to retrieve the Data at URL
// and a convenient interface for persisting the data retrieved
type GGet struct {
	URL      *url.URL
	Strategy strategy.Handler
}

// Execute will execute gget with the provided strategy
func (g *GGet) Execute() (err error) {
	var b []byte

	if b, err = g.File(g.URL); err != nil {
		return err
	}

	return g.Strategy.Handle(b)
}

// File will perform an HTTP GET on the specified URL and return the raw bytes
// If you use this method, you do not need to supply a save.Saver strategy
// and you may handle the file as you wish, use Execute if you wish to utilize a
// prebuilt strategy for persisting your file.
func (g *GGet) File(uri *url.URL) (b []byte, err error) {
	return getFile(uri)
}

func getFile(url *url.URL) (b []byte, err error) {
	var resp *http.Response
	if resp, err = http.Get(url.String()); err != nil {
		return nil, fmt.Errorf("Error fetching %s: %v", url.String(), err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error fetching %s. Error Code was: %d", url.String(), resp.StatusCode)
	}

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(fmt.Errorf("%s: %s", "Critical error reading response body", err.Error()))
	}
	defer resp.Body.Close()

	return b, nil
}