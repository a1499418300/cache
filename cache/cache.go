package cache

import (
	"sync"
	"time"
)

/**
⽀持过期时间和最⼤内存⼤⼩的的内存缓存库。 按照要求实现这⼀个接⼝。
*/
type Cache interface {
	// size 是⼀个字符串。⽀持以下参数: 1KB，100KB，1MB，2MB，1GB 等
	SetMaxMemory(size string) bool
	// 设置⼀个缓存项，并且在expire时间之后过期
	Set(key string, val interface{}, expire ...time.Duration) bool
	// 获取⼀个值
	Get(key string) (interface{}, bool)
	// 删除⼀个值
	Del(key string) bool
	// 检测⼀个值 是否存在
	Exists(key string) bool
	// 清空所有值
	Flush() bool
	// 返回所有的key 多少
	Keys() int64
}

// cacheItem 内存对象封装
type cacheItem struct {
	mu          sync.Mutex  // 互斥锁
	Data        interface{} // 数据
	CreateTime  time.Time   // 创建时间
	ExpireTime  time.Time   // 过期时间
	LastGetTime time.Time   // 最近访问时间
	getCounter  counter     // 总访问次数
}

func (s *cacheItem) isExpire() bool {
	return s.ExpireTime.UnixNano() < time.Now().UnixNano()
}

func (s *cacheItem) writeVisit() {
	s.mu.Lock()
	s.LastGetTime = time.Now()
	s.mu.Unlock()
	s.getCounter.add(1)
}

const maxDuration = time.Duration(3600 * 24 * 365)

func newCacheItem(data interface{}, duration ...time.Duration) *cacheItem {
	dt := maxDuration
	if len(duration) != 0 {
		dt = duration[0]
	}
	tnow := time.Now()
	return &cacheItem{
		mu:          sync.Mutex{},
		Data:        data,
		CreateTime:  tnow,
		ExpireTime:  tnow.Add(dt),
		LastGetTime: tnow,
		getCounter:  counter{},
	}
}

type counter struct {
	sync.RWMutex     // 读写锁
	num          int // 计数器
}

func (s *counter) add(n int) {
	s.Lock()
	defer s.Unlock()
	s.num += n
}

func (s *counter) reset() {
	s.Lock()
	defer s.Unlock()
	s.num = 0
}

func (s *counter) get() int {
	s.RLock()
	defer s.RUnlock()
	return s.num
}

func NewCache() Cache {
	return &SyncMapCache{}
}
