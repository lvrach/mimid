package mimid

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type FakeServer struct {
	resources   []Resource
	resourceLog []string
}

func NewFakeServer(testData string) http.Handler {

	fk := &FakeServer{}
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

func (fk *FakeServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	resource, ok := fk.findResource(req)
	if !ok {
		log.Println(req.URL.Path)
		res.WriteHeader(400)
		return
	}

	log.Println(resource.Name())

	fk.resourceLog = append(fk.resourceLog, resource.Name())

	res.WriteHeader(resource.Status)
	res.Write([]byte(resource.Response))
}

func (fk *FakeServer) findResource(req *http.Request) (Resource, bool) {
	candidates := make([]Resource, 0)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}

	for _, r := range fk.resources {
		if r.Verb == req.Method &&
			r.Path == req.URL.Path &&
			r.Data == string(body) &&
			r.IsAfter(fk.resourceLog) {
			candidates = append(candidates, r)
		}
	}
	if len(candidates) == 0 {
		log.Println(string(body))
		return Resource{}, false
	}

	sort.Sort(bySpecificity(candidates))

	return candidates[0], true
}
