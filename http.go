package leecache

import (
	pb "LeeCache/leecachepb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strings"
)

type HTTPServe struct {
	localhost string
}

func NewHTTPServe(host string) *HTTPServe {
	return &HTTPServe{
		localhost: host,
	}
}

func (s *HTTPServe) Log(format string, v ...interface{}) {
	log.Printf("[HTTPServe %s] %s", s.localhost, fmt.Sprintf(format, v...))
}

func (s *HTTPServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log
	s.Log("%s %s", r.Method, r.URL.Path)

	// path : /
	if r.URL.Path == "/" {
		w.Write([]byte("[LeeCache] HTTP serve running."))
		return
	}

	// path : /<groupName>/<key>
	pathinfo := strings.SplitN(r.URL.Path[1:], "/", 2)
	if len(pathinfo) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := pathinfo[0]
	key := pathinfo[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := proto.Marshal(&pb.GetResponse{Value: view.ByteCopy()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}
