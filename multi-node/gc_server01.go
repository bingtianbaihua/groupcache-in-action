package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/golang/groupcache"
)

/*
   此处设置 peers.

   go run gc_server[01/02/03].go
   curl -v 'localhost:[8001/8002/8003]/gc?key=fakekey' 即可
*/

func main() {
	var peer_addr = []string{"http://127.0.0.1:8001", "http://127.0.0.1:8002", "http://127.0.0.1:8003"}
	peer := groupcache.NewHTTPPool("http://127.0.0.1:8001")
	peer.Set(peer_addr...)

	gc := groupcache.NewGroup("multi-node-groupcache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			log.Printf("get %s from db", key)

			if key != "bad" {
				dest.SetString(key + " : fake value")
				return nil
			}
			return errors.New("illegal key")
		}))

	http.HandleFunc("/gc", func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		k := r.URL.Query().Get("key")
		err := gc.Get(nil, k, groupcache.AllocatingByteSliceSink(&data))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write(data)
	})

	port := ":8001"
	log.Printf("multi node groupcache server-01 run in %v http port", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
