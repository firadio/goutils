package utils

import (
	"strings"
	"sync"
)

// 指定的接口缓存
type CacheConf struct {
	Timeout float64    `json:"timeout"`
	Fields  []string   `json:"fields"`
	NoCache [][]string `json:"noCache"`
}

var mCacheConfByService map[string]CacheConf = make(map[string]CacheConf)

type CacheInfo struct {
	ModifyTime float64
	Data       []byte
}

var mResByService map[string]CacheInfo = make(map[string]CacheInfo)
var mutex sync.Mutex

func InitCache(mCacheConf map[string]CacheConf) {
	mCacheConfByService = mCacheConf
}
func GetResCache(service string, param map[string]interface{}) []byte {
	cacheConf, configed := mCacheConfByService[service]
	if !configed {
		// 缓存未配置
		return nil
	}
	for _, kv := range cacheConf.NoCache {
		if paramValue, ok := param[kv[0]]; ok {
			if paramValue == kv[1] {
				return nil
			}
		}
	}
	mutex.Lock()
	oResInfo, Cached := mResByService[getCacheKey(service, cacheConf.Fields, param)]
	mutex.Unlock()
	if !Cached {
		// 没有缓存过
		return nil
	}
	if TimestampFloat64()-oResInfo.ModifyTime > cacheConf.Timeout {
		// 超时的缓存要丢弃
		return nil
	}
	return oResInfo.Data
}
func PutResCache(service string, param map[string]interface{}, data []byte) {
	cacheConf, configed := mCacheConfByService[service]
	if !configed {
		// 没配置缓存的不保存
		return
	}
	oResInfo := CacheInfo{}
	oResInfo.Data = data
	oResInfo.ModifyTime = TimestampFloat64()
	mutex.Lock()
	mResByService[getCacheKey(service, cacheConf.Fields, param)] = oResInfo
	mutex.Unlock()
}
func getCacheKey(service string, fields []string, param map[string]interface{}) string {
	aCacheKeys := []string{}
	aCacheKeys = append(aCacheKeys, service)
	for _, paramName := range fields {
		if paramValue, ok := param[paramName]; ok {
			// 因为遍历顺序是根据配置项，所以append进来的顺序是固定的
			s := InterfaceToString(paramValue)
			aCacheKeys = append(aCacheKeys, s)
		}
	}
	sCacheKey := strings.Join(aCacheKeys, ",")
	return sCacheKey
}
