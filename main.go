package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// ParseJSON 转换
func ParseJSON(jsonStr string, param interface{}) interface{} {
	json.Unmarshal([]byte(jsonStr), &param)
	fmt.Println("转换结果：", param)
	return param
}

func main() {
	// jsonStr := `{"age": 123, "name": "jack"}`
	// var respMap1 map[string]interface{}
	// respMap1 = ParseJSON(jsonStr, respMap1)
	// fmt.Println("返回结果", respMap1)
	// jsonStr = `[{"age": 123, "name": "jack"}]`
	// var respMap2 []map[string]interface{}
	// respMap2 = ParseJSON(jsonStr, respMap2)
	// fmt.Println("返回结果", respMap2)

	goFunc1(f1)
	goFunc2(f2, 100)

	goFunc(f1)
	goFunc(f2, "xxxx")
	goFunc(f3, "hello", "world", 1, 3.14)
	time.Sleep(5 * time.Second)
}

func goFunc1(f func()) {
	go f()
}

func goFunc2(f func(interface{}), i interface{}) {
	go f(i)
}

func goFunc(f interface{}, args ...interface{}) {
	if len(args) > 1 {
		go f.(func(...interface{}))(args)
	} else if len(args) == 1 {
		go f.(func(interface{}))(args[0])
	} else {
		go f.(func())()
	}
}

func f1() {
	fmt.Println("f1 done")
}

func f2(i interface{}) {
	fmt.Println("f2 done", i)
}

func f3(args ...interface{}) {
	fmt.Println("f3 done", args)
}
