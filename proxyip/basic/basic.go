package proxy

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/firadio/goutils/utils"
)

type Class struct {
	qty            int
	list           []*ProxyInfo
	url_iplist_get string
	Enable         bool
}

func New() *Class {
	proxyip := &Class{}
	return proxyip
}

func (this *Class) SetURL(url_iplist_get string) {
	this.Enable = true
	this.url_iplist_get = url_iplist_get
}

type ProxyInfo struct {
	SocksAddr string
	SocksPort int
}

var Mutex sync.Mutex

func (proxyip *Class) PutList(ipport string) {
	socks5Arr := strings.Split(ipport, ":")
	if len(socks5Arr) != 2 {
		return
	}
	socksPort, err := strconv.Atoi(socks5Arr[1])
	if err != nil {
		return
	}
	socksInfo := &ProxyInfo{SocksAddr: socks5Arr[0], SocksPort: socksPort}
	proxyip.list = append([]*ProxyInfo{socksInfo}, proxyip.list...)
}

func (proxyip *Class) ProxyGetText() string {
	proxyInfo1 := proxyip.ProxyGetOne()
	if proxyInfo1 == nil {
		time.Sleep(time.Second)
		fmt.Println("ProxyGetText获取IP失败，开始重试")
		return proxyip.ProxyGetText()
	}
	return proxyInfo1.SocksAddr + ":" + strconv.Itoa(proxyInfo1.SocksPort)
}

func (proxyip *Class) ProxyGetOne() *ProxyInfo {
	//del_whitelist_by_remark("golang")
	//return
	Mutex.Lock()
	defer Mutex.Unlock()
	if len(proxyip.list) == 0 {
		aLine := proxyip.user_get_ip_list()
		for _, ipaddrport := range aLine {
			//fmt.Println(ipaddrport)
			proxyip.list = append(proxyip.list, ipaddrport)
		}
	}
	if len(proxyip.list) == 0 {
		return nil
	}
	item := proxyip.list[0] // 先进先出
	proxyip.list = proxyip.list[1:len(proxyip.list)]
	return item
}

func (proxyip *Class) user_get_ip_list() []*ProxyInfo {
	aRet := []*ProxyInfo{}
	sUrl := proxyip.url_iplist_get
	_, clientResBody, err := utils.HttpRequestByte("GET", sUrl, nil, nil, nil)
	if err != nil {
		fmt.Println("user_get_ip_list", err)
		return aRet
	}
	sResBody := string(clientResBody)
	sResBody = strings.ReplaceAll(sResBody, "<br/>", "\r\n")
	lines := strings.Split(sResBody, "\r\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		row := strings.Split(line, ":")
		if len(row) != 2 {
			continue
		}
		port, err := strconv.Atoi(row[1])
		if err != nil {
			continue
		}
		socksInfo := &ProxyInfo{SocksAddr: row[0], SocksPort: port}
		aRet = append(aRet, socksInfo)
	}
	return aRet
}
