package storage

import "github.com/wgyuuu/storage_key"
import "github.com/0studio/redisapi"

type ComplexRedisStorage struct {
	RedisStorage
}

func NewComplexRedisStorage(redisClient redisapi.Redis, keyPrefix string, defaultExpireTime int, encoding RedisEncoding) ComplexRedisStorage {
	redisConn := NewRedisStorage(redisClient, keyPrefix, defaultExpireTime, encoding)
	return ComplexRedisStorage{RedisStorage: redisConn}
}

func (this ComplexRedisStorage) GetKeyList(key storage_key.Key) (keyList []storage_key.Key, err error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	keys, err := this.client.Hkeys(cacheKey)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		keyList = append(keyList, storage_key.NewKeyForString(key))
	}
	return
}

func (this ComplexRedisStorage) SetKeyList(key storage_key.Key, keyList []storage_key.Key) error {
	return nil
}
