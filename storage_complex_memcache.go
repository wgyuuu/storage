package storage

import (
	"github.com/dropbox/godropbox/memcache"
	"github.com/wgyuuu/storage"
	"github.com/wgyuuu/storage/keylist"
	"github.com/wgyuuu/storage_key"
)

func (this MemcacheStorage) GetKeyList(key storage_key.Key) ([]storage_key.Key, error) {
	cacheKey, err := storage.BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return []storage_key.Key{}, err
	}

	item := this.client.Get(cacheKey)
	if item.Error() != nil || item.Status() != memcache.StatusNoError {
		return []storage_key.Key{}, item.Error()
	}

	var kl keylist.Keylist
	err := kl.Unmarshal(item.Value())
	if err != nil {
		return []storage_key.Key{}, err
	}

	return storage_key.KeyListForStrList(kl.GetKeyList()), nil
}

func (this MemcacheStorage) SetKeyList(key storage_key.Key, keyList []storage_key.Key) error {
	var kl keylist.Keylist
	for _, key := range keyList {
        kl.KeyList == append(kl.KeyList, key.ToString())
	}

    return this.Set(key, kl)
}
