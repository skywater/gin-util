package main

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/skywater/gin-util/stringUtil"
)

func main() {
	// goFunc1(f1)
	// goFunc2(f2, 100)

	// goFunc(f1)
	// goFunc(f2, "xxxx")
	// goFunc(f3, "hello", "world", 1, 3.14)
	// time.Sleep(5 * time.Second)

	str := `{"max_limit":10537162.340000,"max_type_limit":1000000.000000,"update_time":"2020-06-12 09:59:32","pd_limit":1000000.000000,"create_time":"2020-06-12 09:59:32","city_limit":210743.246800,"credit_amount":120000.000000,"cust_name":"济南泉泽照明电器有限公司","base_limit":363270.559000,"remark":"返回发票参数:{\"enterpriseName\":\"济南泉泽照明电器有限公司\",\"invoiceMonth\":36,\"lst12MonFullAmt\":10537162.3400,\"lst12MonthDecre\":0.6228,\"lst6MonthDecre\":0.0037,\"staticWlMeet12AndZero\":true,\"taxNumber\":\"91370105689848019Q\",\"wlLst10MonAvg\":{\"dynamicWhiteStatistic\":{\"grade1Avg\":0.0,\"grade2Avg\":0.0,\"grade3Avg\":4566.016,\"grade4Avg\":143818.744,\"grade5Avg\":173602.041},\"staticWhiteStatistic\":{\"grade1Avg\":254261.052,\"grade2Avg\":66029.531,\"grade3Avg\":0.0,\"grade4Avg\":0.0,\"grade5Avg\":0.0}}}返回发票PD:78.97返回注册城市编码:{\"companyName\":\"济南泉泽照明电器有限公司\",\"regiCityCode\":370100,\"regiProvinceCode\":370000,\"taxNo\":\"91370105689848019Q\"}"}`
	ss := stringUtil.ParseJSON(str, nil)
	ss1 := *(ss.(*map[string]interface{}))
	remark := ss1["remark"].(string)
	// fmt.Println(remark)
	regexStr := `(.+?)[:：](\{\".+?\":\{.+?\}\}\})|(.+?)[:：](\{\".+?\":.+?\})`
	// regexStr = `\w+_\w+`
	reg := regexp.MustCompile(regexStr)
	// lst := reg.FindAllStringSubmatch(remark, -1)
	// fmt.Println(lst)

	lst := reg.FindStringSubmatch(remark)
	groupNames := reg.SubexpNames()
	fmt.Printf("%v, %d, %d\n", groupNames, len(lst), len(groupNames))
	// fmt.Println(lst)
	// fmt.Println(stringUtil.ToPrettyJSON(match))
	// panic: assignment to entry in nil map，必须先赋值
	// var jsonMap map[string]interface{}
	jsonMap := make(map[string]string)
	for i, v := range lst {
		fmt.Println(i, v)
		if stringUtil.IsNotBlank(v) && i%2 == 0 {
			jsonMap[stringUtil.JoinStr("key_", i)] = v
			respMapJ := stringUtil.ParseJSON(lst[i+1], nil)
			fmt.Println(respMapJ)
			// resp := respMapJ.(*new(map[string]interface{}))
			// jsonMap[stringUtil.JoinStr("val_", i)] = resp
			i++
		}
	}
	fmt.Println(stringUtil.ToPrettyJSON(jsonMap))
}

// RegAllMatchs 循环获取所有匹配
// func RegAllMatchs() string {
//   var index = 0;
//   var ret = [];
//   while(index < data.length && index >= 0){
// 	index = -1;
// 	var match = regexp.MatchString(data);
// 	if(null != match){
// 	  index = match.index;
// 	  data = data.substr(index + match[0].length);
// 	  ret.push(match[0]);
// 	}
//   }
//   return ret
// }

func goFunc1(f func()) {
	ty := reflect.TypeOf(f)
	fmt.Println(f, ty, ty.Kind(), ty.Name())
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
