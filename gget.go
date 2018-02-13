package gget

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// GGet provides a way to retrieve the Data at URL
// and a convenient interface for persisting the data retrieved
type GGet struct {
	URLS     []*url.URL
	Strategy Handler
	Verbose  bool
	r        retriever
}

// Default returns an opinionated version of GGet
// utilizing a retriever that defaults to an HTTP GET
func Default(u []*url.URL, s Handler, verbose bool) *GGet {
	return &GGet{URLS: u, Strategy: s, r: &httpRetriever{}, Verbose: verbose}
}

// retriever provides the interface for Retrieving a file
type retriever interface {
	// get will retrieve the file from the specified URL as a []byte or an error
	get(u *url.URL) (b []byte, err error)
}

type httpRetriever struct{}

type result struct {
	url *url.URL
	b   []byte
	err error
	rtt time.Duration
}

// get will perform an HTTP GET on the specified URL and return the raw bytes
// this is an internal method which GGet relies on to download a file,
// have another retrieval strategy? make a PR :)
// The consumers of GGet shouldn't concern themselves with how it retrieves
// unless they want to dig into the code
func (r *httpRetriever) get(uri *url.URL) (b []byte, err error) {
	return r.getFile(uri)
}

func (r *httpRetriever) getFile(url *url.URL) (b []byte, err error) {
	var resp *http.Response
	if url == nil {
		return nil, fmt.Errorf("empty url, cannot retrieve a ghost resource")
	}
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

// Execute will execute gget with the provided strategy
func (g *GGet) Execute() error {
	resultCh := make(chan result)
	errs := []error{}
	wg := new(sync.WaitGroup)
	wg.Add(len(g.URLS))
	go func(ch <-chan result, wg *sync.WaitGroup) {
		for {
			select {
			case r, ok := <-ch:
				if !ok {
					ch = nil
					break
				}
				if r.err != nil {
					if g.Verbose {
						log.Printf("Error; %v", r.err)
					}
					errs = append(errs, r.err)
				} else {
					if err := g.Strategy.Handle(r); err != nil {
						if g.Verbose {
							log.Printf("Error; %v", err)
						}
						errs = append(errs, err)
					}
				}
				wg.Done()
			}
			if resultCh == nil {
				break
			}
		}
	}(resultCh, wg)

	for _, v := range g.URLS {
		go func(u *url.URL) {
			start := time.Now()
			b, err := g.r.get(u)
			resultCh <- result{b: b, err: err, url: u, rtt: time.Now().Sub(start)}
		}(v)
	}
	wg.Wait()
	close(resultCh)
	if len(errs) > 0 {
		return fmt.Errorf("%d errors encountered during retrieval, or handling, run with Verbose enabled `--v=true` for more detail", len(errs))
	}
	return nil
}
