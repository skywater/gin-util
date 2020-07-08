package httpUtil

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/skywater/gin-util/stringUtil"
)

// BaseResp 通用返回
type BaseResp struct {
	code int
	msg  string
	data string
}

// DoGet http请求
func DoGet(url string) BaseResp {
	return DoRequest("GET", url, "", nil)
}

// DoPost http请求，reqJSON暂未考虑form-data
func DoPost(url string, reqJSON string) BaseResp {
	return DoPostHead(url, reqJSON, nil)
}

// DoPostHead http请求，reqJSON暂未考虑form-data
func DoPostHead(url string, reqJSON string, header map[string]string) BaseResp {
	return DoRequest("POST", url, reqJSON, header)
}

// DoRequest http请求，reqJSON暂未考虑form-data，header暂未处理
func DoRequest(requestType string, remoteURL string, reqJSON string, header map[string]string) BaseResp {
	if stringUtil.IsBlank(requestType) {
		requestType = "GET"
	} else {
		requestType = strings.ToUpper(requestType)
	}
	isPost := strings.EqualFold("POST", requestType)
	log.Printf("http %s请求开始，url=%s，请求参数=%s", requestType, remoteURL, reqJSON)
	client := &http.Client{}
	if isPost { // 跳过证书验证
		// tr := &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// }
		// client = &http.Client{Transport: tr}
	}
	var request io.Reader
	if isPost {
		// bytes.NewBufferString(reqJSON)也可以，注意：POST必须大写！！！！
		request = strings.NewReader(reqJSON)
	}
	uri, e := url.Parse(remoteURL)
	errMsg := dealError(e)
	if stringUtil.IsNotBlank(errMsg) {
		return BaseResp{code: 500, msg: "请求参数异常：" + errMsg}
	}

	req, e := http.NewRequest(requestType, uri.String(), request)
	errMsg = dealError(e)
	if stringUtil.IsNotBlank(errMsg) {
		return BaseResp{code: 500, msg: "请求参数异常：" + errMsg}
	}
	if isPost {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Add("Connection", "keep-alive")
	// if nil != header {
	// 	for i, v := range header {

	// 	}
	// 	req.Header.Set()
	// }
	resp, e := client.Do(req)
	errMsg = dealError(e)
	defer resp.Body.Close()
	if stringUtil.IsNotBlank(errMsg) {
		return BaseResp{code: 500, msg: "请求参数异常：" + errMsg}
	}
	// body, _ := ioutil.ReadAll(resp.Body)
	var body []byte
	if resp.StatusCode == 200 {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip": // 下载文件，暂不处理！！！！！
			reader, _ := gzip.NewReader(resp.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					panic(err)
				}

				if n == 0 {
					break
				}
				body = append(body, buf...)
			}
		default:
			body, _ = ioutil.ReadAll(resp.Body)
		}
	}
	respData := string(body)
	log.Printf("http %s请求结束，url=%s，返回状态=%s，返回参数=%s", requestType, remoteURL, resp.Status, respData)
	baseResp := BaseResp{code: resp.StatusCode, msg: resp.Status, data: respData}
	return baseResp
}
