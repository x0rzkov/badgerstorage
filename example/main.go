package main

import (
    "log"

    "github.com/gocolly/colly"
    "github.com/gocolly/colly/queue"
    "github.com/elonzh/colly-bolt-storage/colly/bolt"
	"go.etcd.io/bbolt"
)

func main() {
    urls := []string{
        "http://httpbin.org/",
        "http://httpbin.org/ip",
        "http://httpbin.org/cookies/set?a=b&c=d",
        "http://httpbin.org/cookies",
    }

    c := colly.NewCollector()
    path := "colly_storage.boltdb"
    var (
        db *bbolt.DB
        err error
    )
    if db, err = bbolt.Open(path, 0666, nil); err != nil {
		panic(err)
	}
    // create the storage
    storage := bolt.NewStorage(db)

    // add storage to the collector
    err = c.SetStorage(storage)
    if err != nil {
        panic(err)
    }

    // close
    defer db.Close()

    // create a new request queue
    q, _ := queue.New(2, storage)

    c.OnResponse(func(r *colly.Response) {
        log.Println("Cookies:", c.Cookies(r.Request.URL.String()))
    })

    // add URLs to the queue
    for _, u := range urls {
        q.AddURL(u)
    }
    // consume requests
    q.Run(c)
}
