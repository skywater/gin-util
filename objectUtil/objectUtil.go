package objectUtil

import (
	"strings"

	"github.com/skywater/gin-util/stringUtil"
)

// LinkedMap 顺序Map；
// 用法，lm := new(LinkedMap); lmObj := lm.NewDefMap()
type LinkedMap struct {
	keys []string               `map的keys`
	mapV map[string]interface{} `map存储`
}

// NewDefMap 初始化空map
func (m *LinkedMap) NewDefMap() LinkedMap {
	return LinkedMap{
		keys: make([]string, 0),
		mapV: make(map[string]interface{}),
	}
}

// NewMap 初始化
func (m *LinkedMap) NewMap(key string, val interface{}) LinkedMap {
	return LinkedMap{
		keys: []string{key},
		mapV: map[string]interface{}{key: val},
	}
}

// Put 放置数据
func (m *LinkedMap) Put(key string, val interface{}) *LinkedMap {
	if _, ok := m.mapV[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.mapV[key] = val
	return m
}

// Get 取出数据
func (m *LinkedMap) Get(key string) interface{} {
	return m.mapV[key]
}

// Del 删除数据
func (m *LinkedMap) Del(key string) *LinkedMap {
	delete(m.mapV, key)
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
		ret[i] = m.mapV[k]
	}
	return ret
}

// ToString 转为string
func (m *LinkedMap) ToString() string {
	if len(m.keys) == 0 {
		return ""
	}
	var b strings.Builder
	length := len(m.keys)
	b.Grow(length * 10)
	b.WriteString("{")
	for i, k := range m.keys {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(stringUtil.ToJSON(m.mapV[k]))
		if i != length-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("}")
	return b.String()
}
