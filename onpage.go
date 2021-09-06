package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

var (
	client *http.Client

	wg sync.WaitGroup
)
var pattern string
var reverse bool

func main() {

	var concurrency int
	flag.BoolVar(&reverse, "r", false, "Match anything other than the supplied pattern")
	flag.IntVar(&concurrency, "c", 50, "Concurrency")
	flag.StringVar(&pattern, "p", "", " Regex (The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages.)")
	flag.Parse()

	var input io.Reader
	input = os.Stdin

	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Printf("failed to open file: %s\n", err)
			os.Exit(1)
		}
		input = file
	}

	sc := bufio.NewScanner(input)

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        concurrency,
			MaxIdleConnsPerHost: concurrency,
			MaxConnsPerHost:     concurrency,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 5 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	semaphore := make(chan bool, concurrency)

	for sc.Scan() {
		raw := sc.Text()
		wg.Add(1)
		semaphore <- true
		go func(raw string) {
			defer wg.Done()
			u, err := url.ParseRequestURI(raw)
			if err != nil {
				return
			}
			fetchURL(u)

		}(raw)
		<-semaphore
	}

	wg.Wait()

	if sc.Err() != nil {
		fmt.Printf("error: %s\n", sc.Err())
	}
}

func fetchURL(u *url.URL) (*http.Response, error) {

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "matcher/0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if respbody, err := ioutil.ReadAll(resp.Body); err == nil {
		bodyString := string(respbody)
		match, _ := regexp.MatchString(pattern, bodyString)
		if !reverse {
			if match {
				fmt.Println(u.String())
			}
		} else {
			if !match {
				fmt.Println(u.String())
			}
		}
	}

	io.Copy(ioutil.Discard, resp.Body)

	return resp, err

}
