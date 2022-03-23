package peers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPGetter : get cache from peers by HTTP, implemented Getter interface
type HTTPGetter struct {
	host string
}

func (g *HTTPGetter) Get(group string, key string) ([]byte, error) {
	url := "http://" + g.host + "/" + group + "/" + key
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body fail: %v", err)
	}

	return bytes, nil
}
