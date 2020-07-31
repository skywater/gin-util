package wechatUtil

import (
	"github.com/skywater/gin-util/httpUtil"
	"github.com/skywater/gin-util/stringUtil"
)

// SendQywxReboot 发送至企业微信机器人
func SendQywxReboot(jsonStr string, rebootKey string) {
	// 注意：企业微信机器人接收json参数格式如下，content是字符串，所以必须转义！！！
	sendData := `{"msgtype": "text", "text": {"content": "{\"dd\": \"dff\"}"}}`
	// 注意：下面 2 种都是失败的，Warning: wrong json format！！！！Java应该是可以的。
	// sendData = `{"msgtype": "text", "text": {"content": "` + jsonStr + `"}}`
	// sendData = "{\"msgtype\": \"text\", \"text\": {\"content\": " + jsonStr + "}}"

	// 正确做法
	sendMap := map[string]interface{}{"msgtype": "text"}
	sendMap["text"] = map[string]interface{}{"content": jsonStr}
	sendData = stringUtil.ToPrettyJSON(sendMap)

	rebootURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + rebootKey
	httpUtil.DoPost(rebootURL, sendData)
}
