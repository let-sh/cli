package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/let-sh/cli/info"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}
	wh := withHeader{Header: make(http.Header), rt: rt}
	wh.Set("Cli-Version", info.Version)
	return wh
}

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}

func GetJsonWithPath(url string, path string) (data gjson.Result, err error) {
	client.Transport = WithHeader(client.Transport)
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

func GetLatestVersion(channel string) (version string, err error) {
	resp, err := http.Get("https://install.let-sh.com/version")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	logrus.WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
		"body":        string(body),
	}).Debugln("get latest version")

	for _, latest := range strings.Split(string(body), "\n") {
		// for each line in version file
		switch channel {
		case "beta":
			if strings.Contains(latest, "beta") {
				return strings.TrimSpace(strings.Split(
					latest, ":")[1]), nil
			}
		case "rc":
			if strings.Contains(latest, "rc") {
				return strings.TrimSpace(strings.Split(
					latest, ":")[1]), nil
			}
		default:
			if strings.Contains(latest, "latest") || strings.Contains(latest, "stable") {
				return strings.TrimSpace(strings.Split(
					latest, ":")[1]), nil
			}
		}
	}
	return "", errors.New("channel not found")
}

func GenerateShortUrl(url string) (shortendUrl string, err error) {
	payload := make(map[string]interface{})
	payload["url"] = url
	payloadBytes, _ := json.Marshal(&payload)
	body := bytes.NewBuffer(payloadBytes)
	resp, err := client.Post("https://api.let-sh.com/j/", "application/json", body)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data := gjson.Get(string(respBody), "data")

	return data.String(), err
}
