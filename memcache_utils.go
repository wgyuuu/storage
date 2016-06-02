package storage

import (
	"time"

	jump "github.com/dgryski/go-jump"
	"github.com/dropbox/godropbox/memcache"
	"github.com/dropbox/godropbox/net2"
	"github.com/wgyuuu/storage_key"
)

type LogError func(error)
type LogInfo func(v ...interface{})
type ShardFunc func(key string, numShard int) (ret int)

type MemcacheConfig struct {
	AddrList             []string `json:"addr,omitempty"` // list of "ip:port"
	MaxActiveConnections int32    `json:"max_active_connections,omitempty"`
	MaxIdleConnections   uint32   `json:"max_idle_connections,omitempty"`
	ReadTimeOutMS        int      `json:"read_timeout_ms,omitempty"`
	WriteTimeOutMS       int      `json:"write_timeout_ms,omitempty"`
}

func GetClient(config MemcacheConfig, logError LogError, logInfo LogInfo) (mc memcache.Client) {
	return GetShardClient(config, logError, logInfo, nil)
}
func GetShardClient(config MemcacheConfig, logError LogError, logInfo LogInfo, shardFunc ShardFunc) (mc memcache.Client) {
	if len(config.AddrList) == 0 {
		panic("could not load mc setting,mcAddrList len 0")
	}

	return getClientFromShardPool(config, logError, logInfo, shardFunc)
}

func getClientFromShardPool(config MemcacheConfig, logError LogError, logInfo LogInfo, shardFunc ShardFunc) (mc memcache.Client) {
	options := net2.ConnectionOptions{
		MaxActiveConnections: config.MaxActiveConnections,
		MaxIdleConnections:   config.MaxIdleConnections,
		ReadTimeout:          time.Duration(config.ReadTimeOutMS) * time.Millisecond,
		WriteTimeout:         time.Duration(config.WriteTimeOutMS) * time.Millisecond,
	}

	if shardFunc == nil {
		shardFunc = func(mcKey string, numShard int) (ret int) {
			if numShard == 0 {
				return -1
			}
			if numShard < 2 {
				return 0
			}
			// https://github.com/renstrom/go-jump-consistent-hash
			// jump 一致性hash 算法
			ret = int(jump.Hash(uint64(storage_key.String(mcKey).ToSum()), numShard))
			// ret = int(crc32.ChecksumIEEE([]byte(key))) % len(mcAddrList)
			return
		}
	}

	manager := NewStaticShardManager(
		config.AddrList,
		logError,
		logInfo,
		shardFunc,
		options)
	mc = memcache.NewShardedClient(manager, false)

	return
}

func NewStaticShardManager(serverAddrs []string, logError LogError, logInfo LogInfo, shardFunc ShardFunc, options net2.ConnectionOptions) memcache.ShardManager {
	// 从dropbox/memcache/static_shard_manager.go copy 来
	// 将其中的log 换成zerogame.info/log
	manager := &memcache.StaticShardManager{}
	manager.Init(
		shardFunc,
		logError,
		logInfo,
		options)

	shardStates := make([]memcache.ShardState, len(serverAddrs), len(serverAddrs))
	for i, addr := range serverAddrs {
		shardStates[i].Address = addr
		shardStates[i].State = memcache.ActiveServer
	}

	manager.UpdateShardStates(shardStates)
	return manager
}
