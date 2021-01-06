package utils

import (
	"io/ioutil"
	//"log"
	"net/http"
	"strings"

	"golang.org/x/net/proxy"
)

func ProxyHttp(socks5ipport string, sUrl string) ([]byte, error) {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", socks5ipport, nil, proxy.Direct)
	if err != nil {
		//fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		return nil, err
		//os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	resp, err := httpClient.Get(sUrl)
	//defer resp.Body.Close()
	if err != nil {
		//log.Println("[httpClient.Get]", err)
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func GetLocationByProxy(ipport string) []string {
	sUrl := "http://myip.ipip.net"
	bResult, err := ProxyHttp(ipport, sUrl)
	s := []string{}
	if err != nil {
		return s
	}
	s = strings.Split(string(bResult), "  ")
	ip := strings.Split(s[0], "：")
	if len(ip) < 2 {
		return s
	}
	s[0] = ip[1]
	location := strings.Split(s[1], "：")
	if len(location) < 2 {
		return s
	}
	s[1] = location[1]
	s[2] = strings.TrimRight(s[2], "\n")
	return s
}
