package main

import (
	leecache "LeeCache"
	"log"
	"net/http"
)

//a simple api serve example

func startApiServe(addr string) {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {

		//log
		log.Printf("[Api serve] %s %s", r.Method, r.URL.Path)

		group := r.URL.Query().Get("group")
		key := r.URL.Query().Get("key")

		g := leecache.GetGroup(group)

		if g == nil {
			http.Error(w, "no such group: "+group, http.StatusNotFound)
			return
		}

		view, err := g.Get(key)
		if err != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteCopy())

	})

	log.Printf("[LeeCache] Api serve is running at %s", addr)
	http.ListenAndServe(addr, nil)
}
