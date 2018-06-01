package mimid

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Proxy struct {
	remoteURL      *url.URL
	localURL       *url.URL
	savePath       string
	resources      []Resource
	lastOccurrence map[string]int
}

func NewProxy(local, remote, savePath string) (*Proxy, error) {
	err := os.MkdirAll(savePath, 0770)
	if err != nil {
		return &Proxy{}, err
	}

	remoteURL, err := url.Parse(remote)
	if err != nil {
		return &Proxy{}, err
	}

	localURL, err := url.Parse(local)
	if err != nil {
		return &Proxy{}, err
	}

	return &Proxy{
		remoteURL,
		localURL,
		savePath,
		[]Resource{},
		make(map[string]int),
	}, nil
}

func (p *Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(p.remoteURL)

	bodyBytes, _ := ioutil.ReadAll(req.Body)
	req.Body.Close() //  must close
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	proxy.ModifyResponse = func(remoteRes *http.Response) error {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		p.capture(remoteRes, req)

		return nil
	}

	proxy.ServeHTTP(res, req)
}

func (p *Proxy) capture(res *http.Response, req *http.Request) {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	resBody = bytes.Replace(resBody, []byte(p.remoteURL.String()), []byte(p.localURL.String()), -1)

	res.Header.Set("Content-Length", fmt.Sprint(len(resBody)))
	res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))

	reqBody, _ := ioutil.ReadAll(req.Body)

	r := Resource{
		Verb:     req.Method,
		Path:     req.URL.Path,
		Data:     string(reqBody),
		Header:   parseHeader(filterHeader(req.Header)),
		Response: string(resBody),
		Status:   res.StatusCode,
	}

	r.After = p.detectAfter(r)

	p.resources = append(p.resources, r)

	p.lastOccurrence[r.expectHash()] = len(p.resources) - 1

	r.Save(p.savePath)
}

func (p *Proxy) detectAfter(current Resource) []string {
	after := make([]string, 0)

	expectHash := current.expectHash()
	start, found := p.lastOccurrence[expectHash]

	if !found {
		return after
	}

	if len(p.resources) == 0 {
		return after
	}

	rr := p.resources[start : len(p.resources)-1]
	for _, r := range rr {
		after = append(after, r.Name())
	}

	return after
}

func parseHeader(src http.Header) (dst map[string]string) {
	dst = make(map[string]string)
	for k, vv := range src {
		for _, v := range vv {
			dst[k] = v
		}
	}
	return dst
}

func filterHeader(hdr http.Header) http.Header {
	hdr.Del("Content-Length")
	hdr.Del("User-Agent")
	hdr.Del("Accept")

	return hdr
}
