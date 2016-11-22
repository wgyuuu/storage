package storage

import (
	"github.com/0studio/redisapi"
	"github.com/wgyuuu/storage_key"
)

type RedisEncoding interface {
	Encoding
	GetKey(obj interface{}) storage_key.Key
}

type RedisStorage struct {
	client            redisapi.Redis
	KeyPrefix         string
	DefaultExpireTime int
	encoding          RedisEncoding
}

func NewRedisClient(serverUrl string) (redisapi.Redis, error) {
	return redisapi.InitDefaultClient(serverUrl)
}

func NewRedisStorage(redisClient redisapi.Redis, keyPrefix string, defaultExpireTime int, encoding RedisEncoding) RedisStorage {
	return RedisStorage{redisClient, keyPrefix, defaultExpireTime, encoding}
}

func (this RedisStorage) Get(key storage_key.Key) (interface{}, error) {
	var err error
	var data []byte
	if keyList, ok := key.(storage_key.KeyList); ok {
		cacheKey, _ := BuildCacheKey(this.KeyPrefix, keyList[0])
		obj, errGet := this.client.Hget(cacheKey, key.ToString())
		err = errGet
		data = obj.([]byte)
	} else {
		cacheKey, _ := BuildCacheKey(this.KeyPrefix, key)
		data, err = this.client.Get(cacheKey)
	}
	if err != nil || data == nil {
		return nil, err
	}

	object, err := this.encoding.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (this RedisStorage) Set(key storage_key.Key, object interface{}) error {
	buf, err := this.encoding.Marshal(object)
	if err != nil {
		return err
	}
	objKey := this.encoding.GetKey(object)
	if keyList, ok := objKey.(storage_key.KeyList); ok {
		keyCache, _ := BuildCacheKey(this.KeyPrefix, keyList[0])
		this.client.Hset(keyCache, objKey.ToString(), buf)
	} else {
		keyCache, _ := BuildCacheKey(this.KeyPrefix, key)
		this.client.Set(keyCache, buf)
	}
	return nil
}

func (this RedisStorage) Add(key storage_key.Key, object interface{}) error {
	return this.Set(key, object)
}

func (this RedisStorage) MultiGet(keys []storage_key.Key) (map[storage_key.Key]interface{}, error) {
	result := make(map[storage_key.Key]interface{}, len(keys))
	for _, key := range keys {
		obj, err := this.Get(key)
		if err != nil {
			return nil, err
		}
		result[key] = obj
	}
	return result, nil
}

func (this RedisStorage) MultiSet(valueMap map[storage_key.Key]interface{}) error {
	for key, value := range valueMap {
		err := this.Set(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this RedisStorage) Delete(key storage_key.Key) (err error) {
	if keyList, ok := key.(storage_key.KeyList); ok {
		cacheKey, _ := BuildCacheKey(this.KeyPrefix, keyList[0])
		err = this.client.Hdel(cacheKey, key.ToString())
	} else {
		cacheKey, _ := BuildCacheKey(this.KeyPrefix, key)
		err = this.client.Delete(cacheKey)
	}
	return
}

func (this RedisStorage) FlushAll() {
	this.client.ClearAll()
}
