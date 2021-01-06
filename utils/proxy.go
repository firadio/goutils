package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Proxy struct {
	qty  int
	list []string
}

func ProxyNew(qty int) *Proxy {
	proxy := &Proxy{}
	proxy.qty = qty
	return proxy
}

func (proxy *Proxy) ProxyGetOne() string {
	//del_whitelist_by_remark("golang")
	//return
	if len(proxy.list) == 0 {
		aLine := user_get_ip_list(proxy.qty)
		for _, ipaddrport := range aLine {
			//fmt.Println(ipaddrport)
			proxy.list = append(proxy.list, ipaddrport)
		}
	}
	item := proxy.list[0] // 先进先出
	proxy.list = proxy.list[1:len(proxy.list)]
	return item
}

//填写Token密钥
var token = "v8b42dmWJEm4KAKb1608965085813"

func user_get_ip_list(qty int) []string {
	aRet := []string{}
	sUrl := "http://list.rola-ip.site:8088/user_get_ip_list"
	query := url.Values{}
	query.Add("token", token)
	query.Add("qty", strconv.Itoa(qty)) //每次获取IP数量
	query.Add("country", "us")          //国家
	query.Add("time", "5")              //时效（分钟）
	query.Add("format", "json")         //返回格式
	query.Add("protocol", "socks5")     //筛选格式
	query.Add("filter", "1")            //是否去重
	statusCode, apiResJson, err := HttpRequestJson("GET", sUrl, query, nil)
	if err != nil {
		fmt.Println("user_get_ip_list", err)
		return aRet
	}
	if apiResJson.Code == 101 {
		fmt.Println("GET_IP_LIST", statusCode, apiResJson.Code, apiResJson.Msg)
		aMsg := strings.SplitN(string(apiResJson.Msg), " ", 2)
		ip := aMsg[0]
		//msg := aMsg[1]
		del_whitelist_by_remark("golang")
		user_add_whitelist(ip, "golang")
		time.Sleep(5 * time.Second)
		return user_get_ip_list(qty)
	}
	aRows := apiResJson.Data.([]interface{})
	for _, _mRow := range aRows {
		mRow := _mRow.(string)
		aRet = append(aRet, mRow)
	}
	return aRet
}

func del_whitelist_by_remark(remark string) {
	aRows := user_get_whitelist()
	for _, _mRow := range aRows {
		mRow := _mRow.(map[string]interface{})
		//fmt.Println(mRow)
		sRemark := mRow["remark"].(string)
		sIp := mRow["ip"].(string)
		if sRemark == remark {
			user_del_whitelist(sIp)
		}
	}
}

func user_get_whitelist() []interface{} {
	//查看白名单
	sUrl := "http://admin.rola-ip.co/user_get_whitelist"
	query := url.Values{}
	query.Add("token", token)
	_, jsonResponse, err := HttpRequestJson("GET", sUrl, query, nil)
	if err != nil {
		fmt.Println("user_get_whitelist", err)
		return nil
	}
	aRows := jsonResponse.Data.([]interface{})
	return aRows
}

func user_add_whitelist(ip string, remark string) {
	//添加白名单
	sUrl := "http://admin.rola-ip.co/user_add_whitelist"
	query := url.Values{}
	query.Add("token", token)
	query.Add("ip", ip)
	query.Add("remark", remark)
	statusCode, jsonResponse, err := HttpRequestJson("GET", sUrl, query, nil)
	if err != nil {
		fmt.Println("user_add_whitelist", err)
		return
	}
	fmt.Println("user_add_whitelist", statusCode, jsonResponse.Code, jsonResponse.Data, jsonResponse.Msg)
}

func user_del_whitelist(ip string) {
	//删除白名单
	sUrl := "http://admin.rola-ip.co/user_del_whitelist"
	query := url.Values{}
	query.Add("token", token)
	query.Add("ip", ip)
	statusCode, jsonResponse, err := HttpRequestJson("GET", sUrl, query, nil)
	if err != nil {
		fmt.Println("user_del_whitelist", err)
		return
	}
	fmt.Println("user_del_whitelist", statusCode, jsonResponse.Code, jsonResponse.Data, jsonResponse.Msg)
}

/*
HTTP处理用函数
*/
type JsonRequest struct {
}

type JsonResponse struct {
	Code int32       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func HttpRequestJson(method string, sUrl string, query url.Values, userReqJson interface{}) (int, JsonResponse, error) {
	//以JSON形式提交并返回JSON
	apiResJson := JsonResponse{}
	bReqData, err := json.Marshal(userReqJson)
	if err != nil {
		return 0, apiResJson, err
	}
	statusCode, clientResBody, err := HttpRequestByte(method, sUrl, query, bytes.NewBuffer(bReqData))
	if err != nil {
		return statusCode, apiResJson, err
	}
	err = json.Unmarshal(clientResBody, &apiResJson)
	if err != nil {
		return statusCode, apiResJson, err
	}
	return statusCode, apiResJson, nil
}

func HttpRequestByte(method string, sUrl string, query url.Values, body io.Reader) (int, []byte, error) {
	//以byte数组返回结果
	if query != nil {
		sUrl += "?" + query.Encode()
	}
	clientReq, err := http.NewRequest(method, sUrl, body)
	if err != nil {
		return 0, nil, err
	}
	httpClient := &http.Client{}
	clientRes, err := httpClient.Do(clientReq) //向后端服务器提交数据
	if err != nil {
		return 0, nil, err
	}
	clientResBody, err := ioutil.ReadAll(clientRes.Body) //取得后端服务器返回的数据
	clientRes.Body.Close()
	if err != nil {
		return clientRes.StatusCode, nil, err
	}
	return clientRes.StatusCode, clientResBody, nil
}
