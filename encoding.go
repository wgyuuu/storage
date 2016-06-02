package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/wgyuuu/storage_key"
)

type Encoding interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte) (interface{}, error)
}

type JsonEncoding struct {
	T reflect.Type
}

func (this JsonEncoding) Marshal(v interface{}) ([]byte, error) {
	buf, err := json.Marshal(v)
	return buf, err
}

func (this JsonEncoding) Unmarshal(data []byte) (interface{}, error) {
	tStruct := reflect.New(this.T)
	dec := json.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(tStruct.Interface())
	if err != nil {
		return err, nil
	}
	return reflect.Indirect(tStruct.Elem()).Interface(), nil
}

type GobEncoding struct {
	T reflect.Type
}

func (this GobEncoding) Marshal(v interface{}) (data []byte, err error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err = enc.Encode(v)
	data = network.Bytes()
	return
}

func (this GobEncoding) Unmarshal(data []byte) (v interface{}, err error) {
	var network bytes.Buffer
	tStruct := reflect.New(this.T)
	network.Write(data)
	dec := gob.NewDecoder(&network)
	err = dec.Decode(tStruct.Interface())
	if err != nil {
		fmt.Println(err)
		return
	}

	v = reflect.Indirect(tStruct.Elem()).Interface()
	return
}

type ByteEncoding struct {
}

func (this ByteEncoding) Marshal(v interface{}) (data []byte, err error) {
	data = v.([]byte)
	return
}

func (this ByteEncoding) Unmarshal(data []byte) (v interface{}, err error) {
	v = data
	return
}

func BuildCacheKey(keyPrefix string, key storage_key.Key) (cacheKey string, err error) {
	if key == nil {
		return "", errors.New("key should not be nil")
	}
	return strings.Join([]string{keyPrefix, key.ToString()}, "_"), nil
}

func GetRawKey(Key string) (rawKey storage_key.String) {
	keys := strings.Split(Key, "_")
	return storage_key.String(keys[len(keys)-1])
}

func InitializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			InitializeStruct(ft.Type, f)
		default:
		}
	}
}
