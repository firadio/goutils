package utils

import (
	"encoding/json"
)

type ProxyConfig struct {
	Data ProxyConfigData `json:"data"`
}
type ProxyConfigData struct {
	GoThreads       int    `json:"go_threads"`
	UrlIplistGet    string `json:"url_iplist_get"`
	UrlWhitelistAdd string `json:"url_whitelist_add"`
	UrlWhitelistGet string `json:"url_whitelist_get"`
	UrlWhitelistDel string `json:"url_whitelist_del"`
}

func GetProxyConfig(sUrl string) (ProxyConfig, error) {
	//以JSON形式提交并返回JSON
	config1 := ProxyConfig{}
	method := "GET"
	_, clientResBody, err := HttpRequestByte(method, sUrl, nil, nil)
	if err != nil {
		return config1, err
	}
	err = json.Unmarshal(clientResBody, &config1)
	if err != nil {
		return config1, err
	}
	return config1, nil
}
