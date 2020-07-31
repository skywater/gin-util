package main

import (
	"fmt"
	"log"
	"time"

	"github.com/skywater/gin-util/stringUtil"
)

func main() {
	time.Sleep(5 * time.Second)
	// panic: interface conversion: interface {} is func(...interface {}) (int, error), not func(...interface {})
	// commonUtil.GoFuncArgs(fmt.Println, "hello", "world")

	// panic: interface conversion: interface {} is func(string) bool, not func(interface {})
	// commonUtil.GoFuncArgs(stringUtil.IsBlank, "hell")
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	log.Println("main begin")
	FuncArgs(stringUtil.IsBlank, "hello")
	FuncArgs(f3, "hello", "world", 1, 3.14)
	log.Println("main ok")
}

func FuncArgs(f interface{}, args ...interface{}) {
	if len(args) > 1 {
		go f.(func(...interface{}))(args)
	} else if len(args) == 1 {
		log.Println("FuncArgs begin")
		// go
		f.(func(string) bool)(args[0].(string))
		log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)
		log.Println("FuncArgs over")
	} else {
		go f.(func())()
	}
}

func f3(args ...interface{}) {
	fmt.Println("f3 done", args)
}
