package sqlx

import (
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

func (conn *OrmDB) Insert(obj orm.OrmStruct) int64 {
	obj.BeforeSave() // 设置值
	sm, values := orm.SchemaValues(obj)

	priIdx := sm.PrimaryIndex()
	if priIdx > 0 {
		values[priIdx] = values[0]
	}

	ret := conn.ExecSql(insertSql(sm), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	ct, err := ret.RowsAffected()
	ErrLog(err)
	return ct
}

func (conn *OrmDB) Delete(obj any) int64 {
	sm := orm.Schema(obj)
	val := sm.PrimaryValue(obj)
	ret := conn.ExecSql(deleteSql(sm), val)
	return parseSqlResult(ret, val, conn, sm)
}

func (conn *OrmDB) Update(obj orm.OrmStruct) int64 {
	obj.BeforeSave()
	sm, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := sm.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	ret := conn.ExecSql(updateSql(sm), values...)
	return parseSqlResult(ret, tVal, conn, sm)
}

// 通过给定的结构体字段更新数据
func (conn *OrmDB) UpdateColumns(obj orm.OrmStruct, columns ...string) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(obj))
	sm := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := updateSqlByColumns(sm, &dstVal, columns)
	ret := conn.ExecSql(upSQL, tValues...)
	return parseSqlResult(ret, tValues[len(tValues)-1], conn, sm)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对应ID值的一行记录
func (conn *OrmDB) QueryID(dest any, id any) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))

	sm := orm.SchemaOfType(dstVal.Type())
	sqlRows := conn.QuerySql(selectSqlForID(sm), id)
	defer CloseSqlRows(sqlRows)

	return scanSqlRowsOne(&dstVal, sqlRows, sm)
}

// 对应ID值的一行记录，支持行记录缓存
func (conn *OrmDB) QueryIDCache(dest any, id any) int64 {
	return queryByIdWithCache(conn, dest, id)
}

// 查询一行记录，查询条件自定义
func (conn *OrmDB) QueryRow(dest any, where string, pms ...any) int64 {
	return conn.QueryRow2(dest, "*", where, pms...)
}

func (conn *OrmDB) QueryRow2(dest any, fields string, where string, pms ...any) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))

	sm := orm.SchemaOfType(dstVal.Type())
	sqlRows := conn.QuerySql(selectSqlForOne(sm, fields, where), pms...)
	defer CloseSqlRows(sqlRows)

	return scanSqlRowsOne(&dstVal, sqlRows, sm)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *OrmDB) QueryRows(dest any, where string, pms ...any) int64 {
	return conn.QueryRows2(dest, "*", where, pms...)
}

func (conn *OrmDB) QueryRows2(dest any, fields string, where string, pms ...any) int64 {
	dSliceTyp, dItemType, isPtr, isKV := checkDestType(dest)

	sm := orm.SchemaOfType(dItemType)
	sqlRows := conn.QuerySql(selectSqlForSome(sm, fields, where), pms...)
	defer CloseSqlRows(sqlRows)

	return scanSqlRowsSlice(dest, sqlRows, sm, dSliceTyp, dItemType, isPtr, isKV)
}

// 高级查询，可以自定义更多参数
func (conn *OrmDB) QueryPet(dest any, pet *SelectPet) int64 {
	dSliceTyp, dItemType, isPtr, isKV := checkDestType(dest)

	sm := orm.SchemaOfType(dItemType)
	if pet.Sql == "" {
		pet.Sql = selectSqlForPet(sm, pet)
	}
	sqlRows := conn.QuerySql(pet.Sql, pet.Prams...)
	defer CloseSqlRows(sqlRows)

	return scanSqlRowsSlice(dest, sqlRows, sm, dSliceTyp, dItemType, isPtr, isKV)
}

func (conn *OrmDB) QueryPetCache(dest any, pet *SelectPetCache) int64 {
	return 0
}
