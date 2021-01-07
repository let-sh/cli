package requests

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func GetJsonWithPath(url string, path string) (data gjson.Result, err error) {
	r, err := client.Get(url)
	if err != nil {
		return data, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	data = gjson.Get(string(body), path)
	return data, err
}
