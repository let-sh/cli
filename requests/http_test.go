package requests

import (
	"fmt"
	"testing"
)

func TestGetLatestVersion(t *testing.T) {
	fmt.Println(GetLatestVersion("rc"))
}

func TestHttpClient(t *testing.T) {
	// resp, err := httpClient.Get("https://api.let-sh.com/query")
	// if err != nil {
	// 	panic(err)

	// }
	// defer resp.Body.Close()
	// s, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf(string(s))
}
