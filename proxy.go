package mimid

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Proxy struct {
	remoteURL *url.URL
	savePath  string
}

func NewProxy(endpoint string, savePath string) (Proxy, error) {

	err := os.MkdirAll(savePath, 0770)
	if err != nil {
		return Proxy{}, err
	}

	remoteURL, err := url.Parse(endpoint)
	if err != nil {
		return Proxy{}, err
	}

	return Proxy{
		remoteURL,
		savePath,
	}, nil
}

func (p Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(p.remoteURL)

	bodyBytes, _ := ioutil.ReadAll(req.Body)
	req.Body.Close() //  must close
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	proxy.ModifyResponse = func(remoteRes *http.Response) error {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		p.observe(remoteRes, req)

		return nil
	}

	proxy.ServeHTTP(res, req)
}

func (p Proxy) observe(res *http.Response, req *http.Request) {
	f, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(f))

	bodyBytes, _ := ioutil.ReadAll(req.Body)

	r := Resource{
		Verb:     req.Method,
		Path:     req.URL.Path,
		Data:     string(bodyBytes),
		Header:   parseHeader(filterHeader(req.Header)),
		Response: string(f),
	}

	r.Save(p.savePath)
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

// func (p Proxy) a_ServeHTTP(res http.ResponseWriter, req *http.Request) {
// 	p.remoteURL.Path = req.URL.Path

// 	remoteReq, err := http.NewRequest(req.Method, p.remoteURL.String(), nil)
// 	remoteReq.Header = req.Header
// 	if err != nil {
// 		log.Println(err)
// 		res.WriteHeader(505)
// 	}

// 	remoteRes, err := p.client.Do(remoteReq)
// 	if err != nil {
// 		log.Println(err)
// 		res.WriteHeader(505)
// 	}

// 	p.observe(remoteRes, req)

// 	header := res.Header()
// 	header.Set("", "")
// 	remoteRes.H

// 	res.WriteHeader(remoteRes.StatusCode)

// 	ioutil.ReadAll(io.TeeReader(remoteRes.Body, res))
// }

// func copyHeader(dst, src http.Header) {
// 	for k, vv := range src {
// 		for _, v := range vv {
// 			dst.Add(k, v)
// 		}
// 	}
// }
