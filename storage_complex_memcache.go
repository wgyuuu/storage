package storage

import (
	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage/keylist"
	"github.com/wgyuuu/storage_key"
)

func (this MemcacheStorage) GetKeyList(key storage_key.Key) ([]storage_key.Key, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return []storage_key.Key{}, err
	}
	item := this.client.Get(cacheKey)
	if item.Error() != nil || item.Status() != memcache.StatusNoError {
		return []storage_key.Key{}, item.Error()
	}

	var kl keylist.Keylist
	err = kl.Unmarshal(item.Value())
	if err != nil {
		return []storage_key.Key{}, err
	}

	var keyList []storage_key.Key
	for _, key := range kl.GetKeyList() {
		keyList = append(keyList, storage_key.NewKeyForString(key))
	}
	return keyList, nil
}

func (this MemcacheStorage) SetKeyList(key storage_key.Key, keyList []storage_key.Key) error {
	var kl keylist.Keylist
	for _, key := range keyList {
		kl.KeyList = append(kl.KeyList, key.ToString())
	}

	buf, err := kl.Marshal()
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
