package main

import (
	"cache/cache"
	"flag"
	"time"

	log "github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer log.Flush()
	c := cache.NewCache()
	c.SetMaxMemory("100MB")
	c.Set("int", 1)
	c.Set("bool", false)
	c.Set("data", map[string]interface{}{"a": 1})
	c.Set("int expire", 1, 10*time.Second)
	c.Get("int")
	c.Del("int")
	c.Get("int expire")
	c.Keys()
	time.Sleep(10 * time.Second)
	c.Get("int expire")
	c.Flush()
	c.Keys()
}
