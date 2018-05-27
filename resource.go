package mimid

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type Resource struct {
	Verb     string
	Path     string
	Header   map[string]string
	Data     string
	Status   string
	Response string
}

func (r Resource) Name() string {
	return r.Verb + strings.Replace(r.Path, "/", "_", -1) + "_" + r.expectHash() + "_" + r.resultHash()
}

func (r Resource) expectHash() string {
	h := md5.New()
	h.Write([]byte(r.Data))
	return hex.EncodeToString(h.Sum(nil))[0:7]
}

func (r Resource) resultHash() string {
	h := md5.New()
	h.Write([]byte(r.Response))
	return hex.EncodeToString(h.Sum(nil))[0:7]
}

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
