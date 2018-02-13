package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/jofo8948/gget"
)

func main() {
	var (
		f   *os.File
		err error
	)

	infile := flag.String("i", "urls.txt", "Use a line-separated file to fetch many files at once.")
	verbose := flag.Bool("v", false, "enables verbose mode")
	flag.Parse()
	if f, err = os.Open(*infile); err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var urls []*url.URL
	for sc.Scan() {
		if u, ok := parseURL(sc.Text()); !ok {
			log.Printf("Couldn't parse URL %s", u)
		} else {
			urls = append(urls, u)
		}
	}
	g := gget.Default(urls, &gget.ToFile{Dst: gget.URLBasedPath}, *verbose)
	if err := g.Execute(); err != nil {
		log.Printf(err.Error())
	}
}

func parseURL(line string) (u *url.URL, b bool) {
	var err error
	if line == "" {
		return nil, false
	}
	if u, err = url.Parse(strings.TrimSpace(line)); err != nil {
		return nil, false
	}
	return u, true
}
