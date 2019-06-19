package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/golang/groupcache"
)

/*
   此处不设置 peers. 单机使用 groupcache.

   go run gc_server.go
   curl -v 'localhost:5180/gc?key=fakekey' 即可
*/

func main() {
	group := groupcache.NewGroup("single-node-groupcache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			// 从 db 里获取源数据
			log.Printf("get %s from db", key)

			if key != "bad" {
				dest.SetString(key + " : fake value")
				return nil
			}
			return errors.New("illegal key")
		}))

	http.HandleFunc("/gc", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		log.Println("get key ", key)

		var data []byte
		err := group.Get(nil, key, groupcache.AllocatingByteSliceSink(&data))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write(data)
	})
	port := ":5180"
	log.Printf("single node groupcache run in %v http port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
