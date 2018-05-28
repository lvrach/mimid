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
	resources []Resource
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
		[]Resource{},
	}, nil
}

func (p Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

func (p Proxy) capture(res *http.Response, req *http.Request) {
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
		Status:   res.StatusCode,
	}

	p.resources = append(p.resources, r)

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
