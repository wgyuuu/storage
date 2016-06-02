package storage

import (
	"reflect"

	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage_key"
)

type MemcacheStorage struct {
	client            memcache.Client
	KeyPrefix         string
	DefaultExpireTime int
	encoding          Encoding
}

func NewMemcStorage(client memcache.Client, keyPrefix string, defaultExpireTime int, encoding Encoding) MemcacheStorage {
	return MemcacheStorage{client, keyPrefix, defaultExpireTime, encoding}
}

func (this MemcacheStorage) Get(key storage_key.Key) (interface{}, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	item := this.client.Get(cacheKey)
	if item.Error() != nil || item.Status() != memcache.StatusNoError {
		return nil, item.Error()
	} else {
		object, err := this.encoding.Unmarshal(item.Value())
		if err != nil {
			return nil, err
		} else {
			return object, nil
		}
	}
}

func (this MemcacheStorage) Set(key storage_key.Key, object interface{}) error {
	if object == nil {
		return nil
	}
	if reflect.TypeOf(object).Kind() == reflect.Slice {
		s := reflect.ValueOf(object)
		if s.IsNil() {
			return nil
		}
	}
	buf, err := this.encoding.Marshal(object)
	if err != nil {
		return err
	}
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	response := this.client.Set(&memcache.Item{Key: keyCache, Value: buf, Expiration: uint32(this.DefaultExpireTime)})
	return response.Error()
}

func (this MemcacheStorage) Add(key storage_key.Key, object interface{}) error {
	return this.Set(key, object)
}

func (this MemcacheStorage) MultiGet(keys []storage_key.Key) (map[storage_key.Key]interface{}, error) {
	keyMap := make(map[storage_key.Key]interface{})
	for _, key := range keys {
		keyMap[key] = nil
	}
	cacheKeys := make([]string, len(keyMap))
	i := 0
	for key, _ := range keyMap {
		cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
		if err != nil {
			return nil, err
		}
		cacheKeys[i] = cacheKey
		i = i + 1
	}
	itemMap := this.client.GetMulti(cacheKeys)
	result := make(map[storage_key.Key]interface{})
	for k, item := range itemMap {
		if len(item.Value()) == 0 {
			continue
		}
		object, err := this.encoding.Unmarshal(item.Value())
		if err != nil {
			continue
		}
		result[GetRawKey(k)] = object
	}
	return result, nil
}

func (this MemcacheStorage) MultiSet(objectMap map[storage_key.Key]interface{}) error {
	items := make([]*memcache.Item, 0, len(objectMap))
	for k, v := range objectMap {
		buf, err := this.encoding.Marshal(v)
		if err != nil {
			return err
		}
		keyCache, err := BuildCacheKey(this.KeyPrefix, k)
		item := &memcache.Item{Key: keyCache, Value: buf, Expiration: uint32(this.DefaultExpireTime)}
		items = append(items, item)
	}
	responses := this.client.SetMulti(items)
	for _, response := range responses {
		if response.Error() != nil {
			return response.Error()
		}
	}
	return nil
}

func (this MemcacheStorage) Delete(key storage_key.Key) error {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	response := this.client.Delete(cacheKey)
	return response.Error()
}

func (this MemcacheStorage) FlushAll() {
	this.client.Flush(uint32(0))
}
