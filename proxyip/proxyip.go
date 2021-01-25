package proxyip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/firadio/goutils/utils"
)

type ProxyIP struct {
	qty               int
	list              []*ProxyInfo
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

type ProxyInfo struct {
	SocksAddr string
	SocksPort int
}

var Mutex sync.Mutex

func (proxyip *ProxyIP) ProxyGetOne() *ProxyInfo {
	//del_whitelist_by_remark("golang")
	//return
	Mutex.Lock()
	if len(proxyip.list) == 0 {
		aLine := proxyip.user_get_ip_list(proxyip.qty)
		for _, ipaddrport := range aLine {
			//fmt.Println(ipaddrport)
			proxyip.list = append(proxyip.list, ipaddrport)
		}
	}
	item := proxyip.list[0] // 先进先出
	proxyip.list = proxyip.list[1:len(proxyip.list)]
	Mutex.Unlock()
	return item
}

func (proxyip *ProxyIP) user_get_ip_list(qty int) []*ProxyInfo {
	aRet := []*ProxyInfo{}
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
		socks5Arr := strings.Split(mRow, ":")
		if len(socks5Arr) != 2 {
			continue
		}
		socksPort, err := strconv.Atoi(socks5Arr[1])
		if err != nil {
			continue
		}
		socksInfo := &ProxyInfo{SocksAddr: socks5Arr[0], SocksPort: socksPort}
		aRet = append(aRet, socksInfo)
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
	statusCode, clientResBody, err := utils.HttpRequestByte(method, sUrl, query, bytes.NewBuffer(bReqData), nil)
	if err != nil {
		return statusCode, apiResJson, err
	}
	err = json.Unmarshal(clientResBody, &apiResJson)
	if err != nil {
		return statusCode, apiResJson, err
	}
	return statusCode, apiResJson, nil
}
