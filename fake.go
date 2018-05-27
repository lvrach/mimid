package mimid

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type FakeServer struct {
	resources []Resource
}

func NewFakeServer(testData string) http.Handler {

	fk := FakeServer{}
	fk.load(testData)
	return fk
}

func (fk *FakeServer) load(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		log.Println(f.Name())

		b, err := ioutil.ReadFile(path + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		r := Resource{}
		err = json.Unmarshal(b, &r)
		if err != nil {
			log.Fatal(err)
		}

		fk.resources = append(fk.resources, r)
	}

}

func (fk FakeServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	resource, ok := fk.findResource(req)
	log.Println(req.URL.Path)
	if !ok {
		res.WriteHeader(400)
		return
	}

	res.WriteHeader(200)
	res.Write([]byte(resource.Response))
}

func (fk FakeServer) findResource(req *http.Request) (Resource, bool) {
	for _, r := range fk.resources {
		if r.Verb == req.Method && r.Path == req.URL.Path {
			return r, true
		}
	}

	return Resource{}, false
}
