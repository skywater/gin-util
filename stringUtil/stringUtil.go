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
	return ToJSONStr(v, false)
}

// ToPrettyJSON 格式化json
func ToPrettyJSON(v interface{}) string {
	return ToJSONStr(v, true)
}

// ToJSONStr 转换为json
func ToJSONStr(v interface{}, isPretty bool) string {
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
	json.Unmarshal([]byte(jsonStr), param)
	return param
}

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
		switch e.(type) {
		case string:
			arrStr[i] = e.(string)
		case int:
			arrStr[i] = strconv.Itoa(e.(int))
		case float32, float64:
			arrStr[i] = strconv.FormatFloat(e.(float64), 'f', 4, 32)
		// case float64:
		// 	arrStr[i] = strconv.FormatFloat(e.(float64), 'f', 4, 64)
		default:
			arrStr[i] = strconv.Itoa(e.(int))
		}
	}

	return strings.Join(arrStr, "")
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
