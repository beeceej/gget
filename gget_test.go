package gget

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"
	"testing"
)

type mockRetriever struct {
	mockData []byte
	mockErr  error
}

var mockURL, _ = url.Parse("http://localhost:8080.com")
var mockErrURL, _ = url.Parse("//x")
var ggetExpectError = GGet{
	URLS:     []*url.URL{mockErrURL},
	Strategy: &ToStdOut{},
	r:        &mockRetriever{mockData: nil, mockErr: errors.New("an error")},
	Verbose:  true,
}
var ggetExpectData = GGet{
	URLS:     []*url.URL{mockURL},
	Strategy: &ToStdOut{},
	r:        &mockRetriever{mockData: []byte{0, 1, 1, 1}},
	Verbose:  true,
}

func (r *mockRetriever) get(u *url.URL) ([]byte, error) {
	if r.mockErr != nil {
		return nil, r.mockErr
	}
	return r.mockData, nil
}

func TestGGet(t *testing.T) {
	testGGetExpectError := func() {
		err := ggetExpectError.Execute()
		if err == nil {
			t.Fail()
		}
	}
	testGGetExpectNoError := func() {
		if err := ggetExpectData.Execute(); err != nil {
			t.Fail()
		}
	}

	testGGetExpectError()
	testGGetExpectNoError()
}

func TestHttpRetriever(t *testing.T) {
	expectedResponse := []byte{0, 1, 1}
	tstServer := start([]byte{0, 1, 1})
	defer tstServer.Close()

	localhostURL, _ := url.Parse("http://localhost:8080/test")
	if response, err := new(httpRetriever).get(localhostURL); err != nil {
		t.Fail()
	} else if !bytes.Equal(response, expectedResponse) {
		t.Fail()
	}
}

func start(expectedResponse []byte) *http.Server {
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write(expectedResponse)
	})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server errord out: %s", err)
		}
	}()

	return server
}
