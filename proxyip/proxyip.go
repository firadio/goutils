package proxyip

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/firadio/goutils/http"
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

func NewLocation(aLocation []string) *ProxyIP {
	this := &ProxyIP{}
	this.SetLocation(aLocation)
	return this
}
func randInt(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min)
}

func (this *ProxyIP) SetLocation(aLocation []string) {
	sUrl := "http://list.rola-ip.site:8088/user_get_ip_list?token=v8b42dmWJEm4KAKb1608965085813&type=4g&qty=1&country={country}&state={state}&city={city}&time=10&format=json&protocol=socks5&filter=1"
	randYesOrNo := randInt(0, 1) == 0
	if randYesOrNo {
		sUrl = "http://list.rola-ip.site:8088/user_get_ip_list?token=v8b42dmWJEm4KAKb1608965085813&qty=1&country={country}&state={state}&city={city}&time=10&format=json&protocol=socks5&filter=1"
	}
	sCountry := aLocation[0]
	sUrl = strings.Replace(sUrl, "{country}", url.QueryEscape(sCountry), 1)
	sState := aLocation[1]
	sState = strings.ReplaceAll(sState, " ", "")
	sUrl = strings.Replace(sUrl, "{state}", url.QueryEscape(sState), 1)
	sCity := aLocation[2]
	sCity = strings.ReplaceAll(sCity, " ", "")
	sUrl = strings.Replace(sUrl, "{city}", url.QueryEscape(sCity), 1)
	url_whitelist_add := "http://admin.rola-ip.co/user_add_whitelist?token=v8b42dmWJEm4KAKb1608965085813&remark={remark}&ip={ip}"
	url_whitelist_get := "http://admin.rola-ip.co/user_get_whitelist?token=v8b42dmWJEm4KAKb1608965085813"
	url_whitelist_del := "http://admin.rola-ip.co/user_del_whitelist?token=v8b42dmWJEm4KAKb1608965085813&ip={ip}"
	this.SetURL(sUrl, url_whitelist_add, url_whitelist_get, url_whitelist_del)
}

func (this *ProxyIP) UrlIplistGet() string {
	return this.url_iplist_get
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

func (proxyip *ProxyIP) ProxyGetOneText() string {
	info := proxyip.ProxyGetOne()
	return info.SocksAddr + ":" + strconv.Itoa(info.SocksPort)
}

func (proxyip *ProxyIP) user_get_ip_list(qty int) []*ProxyInfo {
	sUrl := proxyip.url_iplist_get
	statusCode, apiResJson, err := HttpRequestJson("GET", sUrl, nil, nil)
	if err != nil {
		fmt.Println("user_get_ip_list", err)
		return proxyip.user_get_ip_list(qty)
	}
	if apiResJson.Code == 101 {
		//没有加入白名单
		fmt.Println("GET_IP_LIST", statusCode, apiResJson.Code, apiResJson.Msg)
		aMsg := strings.SplitN(string(apiResJson.Msg), " ", 2)
		ip := aMsg[0]
		//msg := aMsg[1]
		proxyip.del_whitelist_by_remark("golang")
		proxyip.user_add_whitelist(ip, "golang")
		time.Sleep(5 * time.Second)
		return proxyip.user_get_ip_list(qty)
	}
	aRet := []*ProxyInfo{}
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
	http1 := http.New()
	statusCode, clientResBody, err := http1.RequestByte(method, sUrl, query, bReqData, nil)
	if err != nil {
		return statusCode, apiResJson, err
	}
	err = json.Unmarshal(clientResBody, &apiResJson)
	if err != nil {
		return statusCode, apiResJson, err
	}
	return statusCode, apiResJson, nil
}
