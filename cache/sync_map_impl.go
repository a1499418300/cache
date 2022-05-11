package cache

import (
	"cache/utils"
	"sync"
	"time"

	log "github.com/golang/glog"
)

const (
	onceDelKeyNum  = 10 // 内存不足时删除10个key
	delkeyRetryNum = 3  // 内存不足删除key重试次数
)

type SyncMapCache struct {
	maxMemory  int64    // 最大内存
	memCounter counter  // 已用内存计数
	keyCounter counter  // key计数
	data       sync.Map // 数据
}

func (s *SyncMapCache) SetMaxMemory(size string) (ok bool) {
	defer func() {
		log.Infof("SetMaxMemory | size: %v, ok: %v", size, ok)
	}()

	sizeInt, err := utils.ParseMemorySize(size)
	if err != nil {
		log.Errorf("ParseMemorySize failed, err: %v", err)
		return false
	}
	s.maxMemory = sizeInt
	return true
}

func (s *SyncMapCache) Set(key string, val interface{}, expire ...time.Duration) (ok bool) {
	defer func() {
		log.Infof("set | key: %v, val: %v, expire: %v, ok: %v", key, val, expire, ok)
	}()

	item := newCacheItem(val, expire...)
	itemSize, err := utils.GetObjSize(item)
	if err != nil {
		log.Errorf("utils.GetObjSize failed, err: %v", err)
		return false
	}
	// sync.map会冗余一份数据，所以size计算乘2
	itemSize = 2 * itemSize

	counter := 0
	for {
		counter++
		if counter > delkeyRetryNum {
			log.Errorf("delKeys call %v, return", delkeyRetryNum)
			return false
		}
		if int64(itemSize)+int64(s.memCounter.get()) > s.maxMemory {
			// 空间不足则进行删除操作
			s.delKeys()
			continue
		}
		break
	}

	s.memCounter.add(itemSize)
	s.keyCounter.add(1)
	s.data.Store(key, item)
	return true
}

func (s *SyncMapCache) Get(key string) (val interface{}, ok bool) {
	defer func() {
		log.Infof("Get | key: %v, val: %v, ok: %v", key, val, ok)
	}()

	data, ok := s.data.Load(key)
	if !ok {
		return nil, false
	}

	item, ok := data.(*cacheItem)
	if !ok {
		log.Errorf("data.(cacheItem) failed, data: %v", data)
		return nil, false
	}

	if item.isExpire() {
		s.Del(key)
		return nil, false
	}

	item.writeVisit()
	return item.Data, true
}

func (s *SyncMapCache) Del(key string) (ok bool) {
	defer func() {
		log.Infof("Del | key: %v, ok: %v", key, ok)
	}()

	data, ok := s.data.Load(key)
	if !ok {
		return false
	}
	itemSize, _ := utils.GetObjSize(data.(*cacheItem))
	s.data.Delete(key)
	s.keyCounter.add(-1)
	s.memCounter.add(-itemSize)
	return true
}

func (s *SyncMapCache) Exists(key string) (ok bool) {
	defer func() {
		log.Infof("Exists | key: %v, ok: %v", key, ok)
	}()

	_, ok = s.data.Load(key)
	return ok
}

func (s *SyncMapCache) Flush() (ok bool) {
	defer func() {
		log.Infof("Flush | ok: %v", ok)
	}()

	s.data = sync.Map{}
	s.memCounter.reset()
	s.keyCounter.reset()
	return true
}

func (s *SyncMapCache) Keys() (n int64) {
	defer func() {
		log.Infof("Keys | n: %v", n)
	}()
	return int64(s.keyCounter.get())
}

// 淘汰部分数据
func (s *SyncMapCache) delKeys() {
	// 淘汰过期数据
	counter := 0
	s.data.Range(func(key, value interface{}) bool { // return true-继续遍历 false-结束遍历
		item, ok := value.(*cacheItem)
		if !ok {
			return true
		}
		if item.isExpire() {
			counter++
			s.Del(key.(string))
		}
		return counter <= onceDelKeyNum
	})
	if counter > onceDelKeyNum {
		return
	}

	// 淘汰访问次数<=n的数据
	n := 0
	s.data.Range(func(key, value interface{}) bool {
		item, ok := value.(*cacheItem)
		if !ok {
			return true
		}
		if item.getCounter.get() <= n {
			counter++
			s.Del(key.(string))
		}
		return counter <= onceDelKeyNum
	})

	// 其它策略
}
