package storage

import "github.com/wgyuuu/storage_key"

type Storage interface {
	Get(key storage_key.Key) (interface{}, error)
	// 兼容mysql (这个标示可以放在object里)
	Add(key storage_key.Key, object interface{}) error
	Set(key storage_key.Key, object interface{}) error
	MultiGet(keys []storage_key.Key) (map[storage_key.Key]interface{}, error)
	MultiSet(map[storage_key.Key]interface{}) error
	Delete(key storage_key.Key) error
	FlushAll()
}

type StorageProxy struct {
	PreferedStorage Storage
	BackupStorage   Storage
}

func NewStorageProxy(prefered, backup Storage) StorageProxy {
	return StorageProxy{
		PreferedStorage: prefered,
		BackupStorage:   backup,
	}
}

func (this StorageProxy) Get(key storage_key.Key) (interface{}, error) {
	object, err := this.PreferedStorage.Get(key)
	if err != nil {
		object, err = this.BackupStorage.Get(key)
		if err == nil {
			return object, nil
		} else {
			return nil, err
		}
	}
	if object == nil {
		object, err = this.BackupStorage.Get(key)
		if err != nil {
			return nil, err
		}
		if object != nil {
			this.PreferedStorage.Set(key, object)
		}
	}
	return object, nil
}

func (this StorageProxy) Add(key storage_key.Key, object interface{}) error {
	if object != nil {
		err := this.PreferedStorage.Add(key, object)
		if err != nil {
			return err
		}
		err = this.BackupStorage.Add(key, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this StorageProxy) Set(key storage_key.Key, object interface{}) error {
	if object != nil {
		err := this.PreferedStorage.Set(key, object)
		if err != nil {
			return err
		}
		err = this.BackupStorage.Set(key, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this StorageProxy) MultiGet(keys []storage_key.Key) (map[storage_key.Key]interface{}, error) {
	resultMap, err := this.PreferedStorage.MultiGet(keys)
	if err != nil {
		return nil, err
	}
	missedKeyCount := 0
	for _, key := range keys {
		if _, find := resultMap[key]; !find {
			missedKeyCount++
		}
	}
	if missedKeyCount > 0 {
		missedKeys := make([]storage_key.Key, missedKeyCount)
		i := 0
		for _, key := range keys {
			if _, find := resultMap[key]; !find {
				missedKeys[i] = key
				i++
			}
		}
		missedMap, err := this.BackupStorage.MultiGet(missedKeys)
		if err != nil {
			return nil, err
		}
		this.PreferedStorage.MultiSet(missedMap)
		for k, v := range missedMap {
			resultMap[k] = v
		}
	}
	return resultMap, nil
}

func (this StorageProxy) MultiSet(objectMap map[storage_key.Key]interface{}) error {
	err := this.PreferedStorage.MultiSet(objectMap)
	if err != nil {
		return err
	}
	err = this.BackupStorage.MultiSet(objectMap)
	if err != nil {
		return err
	}
	return nil
}

func (this StorageProxy) Delete(key storage_key.Key) error {
	err := this.BackupStorage.Delete(key)
	if err != nil {
		return err
	}
	err = this.PreferedStorage.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (this StorageProxy) FlushAll() {
	this.PreferedStorage.FlushAll()
	this.BackupStorage.FlushAll()
}
