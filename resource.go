package mimid

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Resource struct {
	After    []string
	Verb     string
	Path     string
	Header   map[string]string
	Data     string
	Status   int
	Response string
}

func (r Resource) Name() string {
	return r.Verb + strings.Replace(r.Path, "/", "_", -1) + "_" + r.expectHash() + "_" + r.resultHash()
}

func (r Resource) expectHash() string {
	h := md5.New()
	h.Write([]byte(r.Verb))
	h.Write([]byte(r.Path))
	h.Write([]byte(r.Data))
	return hex.EncodeToString(h.Sum(nil))[0:7]
}

func (r Resource) resultHash() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprint(r.Status)))
	h.Write([]byte(r.Response))
	return hex.EncodeToString(h.Sum(nil))[0:7]
}

// IsAfter determines whether the Resource satisfy after criteria
func (r Resource) IsAfter(calls []string) bool {
	if len(r.After) == 0 {
		return true
	}

	if len(r.After) > len(calls) {
		return false
	}

	lastCalls := calls[len(calls)-len(r.After) : len(calls)-1]

	for i := range lastCalls {
		if strings.Compare(lastCalls[i], r.After[i]) == 0 {
			return false
		}
	}
	return true
}

// Save this resource to its file instead the given path
func (r Resource) Save(path string) {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(path+"/"+r.Name()+".json", data, 0666)
	if err != nil {
		log.Fatalln(err)
	}
}

type bySpecificity []Resource

func (rr bySpecificity) Len() int {
	return len(rr)
}
func (rr bySpecificity) Swap(i, j int) {
	rr[i], rr[j] = rr[j], rr[i]
}
func (rr bySpecificity) Less(i, j int) bool {
	return len(rr[i].After) > len(rr[j].After)
}
