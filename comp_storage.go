package storage

import "github.com/wgyuuu/storage_key"

type CompStorage interface {
	Storage
	GetKeyList(key storage_key.Key) ([]storage_key.Key, error)
	SetKeyList(key storage_key.Key, keyList []storage_key.Key) error
}

type CompStorageProxy struct {
	StorageProxy
}

func NewCompStorageProxy(prefered, backup CompStorage) CompStorageProxy {
	return CompStorageProxy{
		PreferedStorage: prefered,
		BackupStorage:   backup,
	}
}

func (this CompStorageProxy) GetKeyList(key storage_key.Key) ([]storage_key.Key, error) {
	keyList, err := this.PreferedStorage.(CompStorage).GetKeyList(key)
	if err != nil {
		return []storage_key.Key{}, err
	}
	if len(keyList) == 0 {
		keyList, err = this.BackupStorage.(CompStorage).GetKeyList(key)
		if err != nil {
			return []storage_key.Key{}, err
		}
		if len(keyList) > 0 {
			this.PreferedStorage.(CompStorage).SetKeyList(key, keyList)
		}
	}
    return keyList, nil
}

func (this CompStorageProxy) SetKeyList(key storage_key.Key, keyList []storage_key.Key) error {
    err := this.BackupStorage.(CompStorage).SetKeyList(key, keyList)
    if err != nil {
        return err
    }
    return this.PreferedStorage.(CompStorage).SetKeyList(key, keyList)
}
