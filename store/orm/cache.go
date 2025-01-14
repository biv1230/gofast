package orm

import (
	"fmt"
	"github.com/qinchende/gofast/skill/hash"
	"reflect"
	"sync"
	"time"
)

func (ms *ModelSchema) TableName() string {
	return ms.attrs.TableName
}

func (ms *ModelSchema) CacheAll() bool {
	return ms.attrs.CacheAll
}

func (ms *ModelSchema) CachePreFix() string {
	return ms.attrs.cacheKeyFmt
}

func (ms *ModelSchema) CacheLineKey(dbName, id any) string {
	return fmt.Sprintf(ms.attrs.cacheKeyFmt, dbName, id)
}

func (ms *ModelSchema) CacheSqlKey(sql string) string {
	return "Gf#Pet#" + hash.Md5HexString(sql)
}

func (ms *ModelSchema) ExpireS() uint32 {
	return ms.attrs.ExpireS
}

// 可以考虑加上随机 5% 左右的偏差，防止将来缓存统一过期导致缓存雪崩
func (ms *ModelSchema) ExpireDuration() time.Duration {
	return time.Duration(ms.attrs.ExpireS) * time.Second
}

func (ms *ModelSchema) FieldsKV() map[string]int8 {
	return ms.fieldsKV
}

func (ms *ModelSchema) ColumnsKV() map[string]int8 {
	return ms.columnsKV
}

func (ms *ModelSchema) Columns() []string {
	return ms.columns
}

func (ms *ModelSchema) UpdatedIndex() int8 {
	return ms.updatedIndex
}

func (ms *ModelSchema) PrimaryIndex() int8 {
	return ms.primaryIndex
}

func (ms *ModelSchema) AutoIndex() int8 {
	return ms.autoIndex
}

func (ms *ModelSchema) InsertSQL(fn func(*ModelSchema) string) string {
	if ms.insertSQL == "" {
		ms.insertSQL = fn(ms)
	}
	return ms.insertSQL
}

func (ms *ModelSchema) UpdateSQL(fn func(*ModelSchema) string) string {
	if ms.updateSQL == "" {
		ms.updateSQL = fn(ms)
	}
	return ms.updateSQL
}

func (ms *ModelSchema) SelectSQL(fn func(*ModelSchema) string) string {
	if ms.selectSQL == "" {
		ms.selectSQL = fn(ms)
	}
	return ms.selectSQL
}

func (ms *ModelSchema) DeleteSQL(fn func(*ModelSchema) string) string {
	if ms.deleteSQL == "" {
		ms.deleteSQL = fn(ms)
	}
	return ms.deleteSQL
}

func (ms *ModelSchema) PrimaryValue(obj any) any {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	return rVal.FieldByIndex(ms.fieldsIndex[ms.primaryIndex]).Interface()
}

func (ms *ModelSchema) ValueByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Interface()
}

func (ms *ModelSchema) AddrByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Addr().Interface()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 缓存数据表的Schema
var cachedSchemas sync.Map

func cacheSetSchema(typ reflect.Type, val *ModelSchema) {
	cachedSchemas.Store(typ, val)
}

func cacheGetSchema(typ reflect.Type) *ModelSchema {
	if ret, ok := cachedSchemas.Load(typ); ok {
		return ret.(*ModelSchema)
	}
	return nil
}
