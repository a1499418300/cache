package utils

import (
	"bytes"
	"cache/errors"
	"encoding/gob"
	"regexp"
	"strconv"
	"strings"
	"sync"

	log "github.com/golang/glog"
)

type ByteSize int64

const (
	Bit ByteSize = 1 << (10 * iota)
	KB
	MB
	GB
)

var (
	gobBuf         *bytes.Buffer
	gobEncoder     *gob.Encoder
	gobEncoderOnce sync.Once
	memorySizeMap  = map[string]int64{
		"K": int64(KB),
		"M": int64(MB),
		"G": int64(GB),
	}
)

// ParseMemorySize 解析内存单位
func ParseMemorySize(size string) (int64, error) {
	// 兼容大小写
	size = strings.ToUpper(size)
	// 兼容1M(1MB)
	sizeRe, err := regexp.Compile(`^(\d+)([GMK])B?$`)
	if err != nil {
		return 0, err
	}
	if !sizeRe.MatchString(size) {
		return 0, errors.ErrMemorySize
	}
	tmp := sizeRe.FindStringSubmatch(size)
	if len(tmp) < 3 {
		return 0, errors.ErrMemorySize
	}
	sizeInt, err := strconv.Atoi(tmp[1])
	if err != nil {
		return 0, err
	}
	return int64(sizeInt) * memorySizeMap[tmp[2]], nil
}

// GetObjSize 获取对象内存大小
// 参考: https://qa.1r1g.com/sf/ask/3098026571/#44258164
// 编码会损耗性能，另外编码后的大小，和实际大小存在偏差。
// 由于unsafe.SizeOf()只能返回 函数值传递类型的大小，遂采用此写法
func GetObjSize(obj interface{}) (int, error) {
	gobEncoderOnce.Do(newGobEncoder)
	gobBuf.Reset()
	if err := gobEncoder.Encode(obj); err != nil {
		log.Errorf("gobEncoder.Encode failed, err: %v, obj: %v", err, obj)
		return 0, errors.ErrGobEncodeFailed
	}
	return gobBuf.Len(), nil
}

func newGobEncoder() {
	gobBuf = new(bytes.Buffer)
	gobEncoder = gob.NewEncoder(gobBuf)
	gob.Register(map[string]interface{}{})
}
