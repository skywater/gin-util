package stringUtil

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// IsBlank 字符串是否为空
func IsBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// IsNotBlank 字符串是否不为空
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

// ToJSON 转为json
func ToJSON(v interface{}) string {
	return toJSONStr(v, false)
}

// ToPrettyJSON 格式化json
func ToPrettyJSON(v interface{}) string {
	return toJSONStr(v, true)
}

// toJSONStr 转换为json
func toJSONStr(v interface{}, isPretty bool) string {
	var jsonMsg []byte
	var e error
	if isPretty {
		jsonMsg, e = json.MarshalIndent(v, "", "    ")
	} else {
		jsonMsg, e = json.Marshal(v)
	}
	if e != nil {
		log.Println(e)
	}
	return string(jsonMsg)
}

// ParseJSON 转换为对象，param为结构地址，如：&Person；返回结果为地址，须用指针
// 1、结构参数：
// baseRespJs := ParseJSON(jsonStr, new(httpUtil.BaseResp))
// resp := baseRespJs.(*httpUtil.BaseResp)
// fmt.Println("返回结果", resp, resp.Code)
// 2、map参数：
// respMapJ := ParseJSON(jsonStr, new(map[string]interface{}))
// resp := *(respMapJ.(*map[string]interface{}))
// fmt.Println("返回结果", resp, resp["code"])
// 3、数组参数：
// respArrayJ := ParseJSON(jsonStr, new([]map[string]interface{}))
// resp := *(respArrayJ.(*[]map[string]interface{}))
// fmt.Println("返回结果", resp, resp[0]["code"])
func ParseJSON(jsonStr string, param interface{}) interface{} {
	if IsBlank(jsonStr) {
		return nil
	}
	jsonStr = strings.TrimSpace(jsonStr)
	if nil == param {
		if IsArray(jsonStr) {
			param = new([]map[string]interface{})
		} else {
			param = new(map[string]interface{})
		}
	}
	json.Unmarshal([]byte(jsonStr), param)
	fmt.Println(param, &param, GetType(param))
	return param
}

// IsArray 字符串是否是数组型json
func IsArray(param string) bool {
	if IsBlank(param) {
		return false
	}
	param = strings.TrimSpace(param)
	if param[0:1] == "[" {
		return true
	}
	return false
}

// func ParseJSON(jsonStr string, param *interface{}) interface{} {
// 	if IsBlank(jsonStr) {
// 		return nil
// 	}
// 	if nil == param {
// 		if jsonStr[0:1] == "[" {
// 			parType := new([]map[string]interface{})
// 			json.Unmarshal([]byte(jsonStr), parType)
// 			return *parType
// 		}
// 		parType := new(map[string]interface{})
// 		json.Unmarshal([]byte(jsonStr), parType)
// 		return *parType
// 	}
// 	json.Unmarshal([]byte(jsonStr), param)
// 	return *param
// }

// JoinStr 字符串拼接
func JoinStr(strs ...interface{}) string {
	if nil == strs || len(strs) == 0 {
		return ""
	}
	arrStr := make([]string, len(strs))
	// n := 0
	for i, e := range strs {
		ty := reflect.TypeOf(e)
		fmt.Println(e, ty, ty.Kind(), ty.Name())
		// ty.(type)只能在switch使用
		arrStr[i] = ToStr(e)
	}

	return strings.Join(arrStr, "")
}

// ToStr 转为字符串
func ToStr(e interface{}) string {
	switch e.(type) {
	case string:
		return e.(string)
	case int:
		return strconv.Itoa(e.(int))
	case float32, float64:
		return strconv.FormatFloat(e.(float64), 'f', 4, 32)
	default:
		return ToJSON(e)
	}
}

// ToInt 字符串转为整型
func ToInt(str string) int {
	// 空字符串字段、空格都会转为 0
	i, _ := strconv.Atoi(str)
	return i
}

// GetType 获取类型
func GetType(e interface{}) string {
	ty := reflect.TypeOf(e)
	fmt.Println(e, "类型信息：", ty, ty.Kind(), ty.Name())
	return ty.Name()
}

func DealError(e error) string {
	if nil != e {
		log.Println(e)
		return e.Error()
	}
	return ""
}
