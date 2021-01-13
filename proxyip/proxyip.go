package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProxyIP struct {
	qty               int
	list              []string
	url_iplist_get    string
	url_whitelist_add string
	url_whitelist_get string
	url_whitelist_del string
}

func ProxyNew() *ProxyIP {
	proxyip := &ProxyIP{}
	return proxyip
}

func (proxyip *ProxyIP) SetURL(url_iplist_get string, url_whitelist_add string, url_whitelist_get string, url_whitelist_del string) {
	proxyip.url_iplist_get = url_iplist_get
	proxyip.url_whitelist_add = url_whitelist_add
	proxyip.url_whitelist_get = url_whitelist_get
	proxyip.url_whitelist_del = url_whitelist_del
}

func (proxyip *ProxyIP) ProxyGetOne() string {
	//del_whitelist_by_remark("golang")
	//return
	if len(proxyip.list) == 0 {
		aLine := proxyip.user_get_ip_list(proxyip.qty)
		for _, ipaddrport := range aLine {
			//fmt.Println(ipaddrport)
			proxyip.list = append(proxyip.list, ipaddrport)
		}
	}
	item := proxyip.list[0] // 先进先出
	proxyip.list = proxyip.list[1:len(proxyip.list)]
	return item
}

func (proxyip *ProxyIP) user_get_ip_list(qty int) []string {
	aRet := []string{}
	sUrl := proxyip.url_iplist_get
	statusCode, apiResJson, err := HttpRequestJson("GET", sUrl, nil, nil)
	if err != nil {
		fmt.Println("user_get_ip_list", err)
		return aRet
	}
	if apiResJson.Code == 101 {
		fmt.Println("GET_IP_LIST", statusCode, apiResJson.Code, apiResJson.Msg)
		aMsg := strings.SplitN(string(apiResJson.Msg), " ", 2)
		ip := aMsg[0]
		//msg := aMsg[1]
		proxyip.del_whitelist_by_remark("golang")
		proxyip.user_add_whitelist(ip, "golang")
		time.Sleep(5 * time.Second)
		return proxyip.user_get_ip_list(qty)
	}
	aRows := apiResJson.Data.([]interface{})
	for _, _mRow := range aRows {
		mRow := _mRow.(string)
		aRet = append(aRet, mRow)
	}
	return aRet
}

func (proxyip *ProxyIP) del_whitelist_by_remark(remark string) {
	aRows := proxyip.user_get_whitelist()
	for _, _mRow := range aRows {
		mRow := _mRow.(map[string]interface{})
		//fmt.Println(mRow)
		sRemark := mRow["remark"].(string)
		sIp := mRow["ip"].(string)
		if sRemark == remark {
			proxyip.user_del_whitelist(sIp)
		}
	}
}

func (proxyip *ProxyIP) user_get_whitelist() []interface{} {
	//查看白名单
	sUrl := proxyip.url_whitelist_get
	_, jsonResponse, err := HttpRequestJson("GET", sUrl, nil, nil)
	if err != nil {
		fmt.Println("user_get_whitelist", err)
		return nil
	}
	aRows := jsonResponse.Data.([]interface{})
	return aRows
}

func (proxyip *ProxyIP) user_add_whitelist(ip string, remark string) {
	//添加白名单
	sUrl := proxyip.url_whitelist_add
	sUrl = strings.Replace(sUrl, "{ip}", ip, 1)
	sUrl = strings.Replace(sUrl, "{remark}", remark, 1)
	statusCode, jsonResponse, err := HttpRequestJson("GET", sUrl, nil, nil)
	if err != nil {
		fmt.Println("user_add_whitelist", err)
		return
	}
	fmt.Println("user_add_whitelist", statusCode, jsonResponse.Code, jsonResponse.Data, jsonResponse.Msg)
}

func (proxyip *ProxyIP) user_del_whitelist(ip string) {
	//删除白名单
	sUrl := proxyip.url_whitelist_del
	sUrl = strings.Replace(sUrl, "{ip}", ip, 1)
	statusCode, jsonResponse, err := HttpRequestJson("GET", sUrl, nil, nil)
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
