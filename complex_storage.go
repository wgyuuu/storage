package storage

import "github.com/wgyuuu/storage_key"

// 复合主键
type ComplexStorage interface {
	Storage
	GetKeyList(key storage_key.Key) ([]storage_key.Key, error)
	SetKeyList(key storage_key.Key, keyList []storage_key.Key) error
}

type ComplexStorageProxy struct {
	StorageProxy
}

func NewComplexStorageProxy(prefered, backup ComplexStorage) ComplexStorageProxy {
	return ComplexStorageProxy{
		PreferedStorage: prefered,
		BackupStorage:   backup,
	}
}

func (this ComplexStorageProxy) GetKeyList(key storage_key.Key) (keyList []storage_key.Key, err error) {
	keyList, err = this.PreferedStorage.(ComplexStorage).GetKeyList(key)
	if err != nil {
		return
	}
	if len(keyList) > 0 {
		return
	}
	keyList, err = this.BackupStorage.(ComplexStorage).GetKeyList(key)
	if len(keyList) > 0 {
		this.PreferedStorage.(ComplexStorage).SetKeyList(key, keyList)
	}
	return
}

// 不支持插入
func (this ComplexStorageProxy) SetKeyList(key storage_key.Key, keyList storage_key.Key) error {
	err := this.BackupStorage.(ComplexStorage).SetKeyList(key, keyList)
	if err != nil {
		return err
	}
	return this.PreferedStorage.(ComplexStorage).SetKeyList(key, keyList)
}

func (this ComplexStorageProxy) PushKey(key storage_key.Key, inKey storage_key.Key) error {
	keyList, err := this.GetKeyList(key)
	if err != nil {
		return err
	}

	keyList = storage_key.AppendKey(keyList, inKey)
	return this.SetKeyList(key, keyList)
}
