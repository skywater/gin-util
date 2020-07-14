package objectUtil

import (
	"fmt"
	"reflect"
	"strings"

	"gin-util/stringUtil"
)

// LinkedMap 顺序Map；
// 用法，lm := new(LinkedMap); lmObj := lm.NewDefMap()
type LinkedMap struct {
	keys []string
	MapV map[string]interface{} `json:"linkedMap"`
}

// NewLinkedMap 初始化空map
func NewLinkedMap() LinkedMap {
	return LinkedMap{
		keys: make([]string, 0),
		MapV: make(map[string]interface{})}
}

// InitLinkedMap 初始化
func InitLinkedMap(key string, val interface{}) LinkedMap {
	return LinkedMap{
		keys: []string{key},
		MapV: map[string]interface{}{key: val}}
}

// Put 放置数据
func (m *LinkedMap) Put(key string, val interface{}) *LinkedMap {
	if _, ok := m.MapV[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.MapV[key] = val
	return m
}

// Get 取出数据
func (m *LinkedMap) Get(key string) interface{} {
	return m.MapV[key]
}

// Del 删除数据
func (m *LinkedMap) Del(key string) *LinkedMap {
	delete(m.MapV, key)
	for i, v := range m.keys {
		if v == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
	return m
}

// Keys 获取所有 Keys
func (m *LinkedMap) Keys() []string {
	return m.keys
}

// Values 获取所有 Values
func (m *LinkedMap) Values() []interface{} {
	ret := make([]interface{}, len(m.keys))
	for i, k := range m.keys {
		ret[i] = m.MapV[k]
	}
	return ret
}

// ToString 转为string
func (m LinkedMap) String() string {
	if len(m.keys) == 0 {
		return ""
	}
	var b strings.Builder
	length := len(m.keys)
	b.Grow(length * 10)
	b.WriteString("{")
	var val interface{}
	for i, k := range m.keys {
		b.WriteString(k)
		b.WriteString(":")
		val = m.MapV[k]
		switch val.(type) {
		case LinkedMap:
			b.WriteString(val.(LinkedMap).String())
		default:
			b.WriteString(stringUtil.ToStr(m.MapV[k]))
		}
		if i != length-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("}")
	return b.String()
}

// IsArray 判断是否是数组，并返回
func IsArray(v interface{}) (bool, []interface{}) {
	if nil == v {
		return false, nil
	}
	// v.(type) 无法判断是否是数组
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice, reflect.Array:
		var nv []interface{}
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i)
			nv = append(nv, elem.Interface())
		}
		return true, nv
		// 报错：interface conversion: interface {} is []map[string]interface {}, not []interface {}
		// return true, v.([]interface{})
	default:
		return false, nil
	}
}

// ArrayIsEmpty 数组为空
func ArrayIsEmpty(v []interface{}) bool {
	return nil == v || len(v) == 0
}

// DealArrayObj 处理数组
func DealArrayObj(v ...interface{}) []interface{} {
	if ArrayIsEmpty(v) {
		return nil
	}
	nv := make([]interface{}, 0, len(v))
	for _, vv := range v {
		if bl, va := IsArray(vv); bl {
			nv = append(nv, va...)
		} else {
			nv = append(nv, vv)
		}
	}
	return nv
}
