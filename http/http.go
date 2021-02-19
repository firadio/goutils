package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type Class struct {
	HttpClient         *http.Client
	RequestMethod      string
	RequestUrl         string
	RequestHeader      map[string]string
	RequestBody        io.Reader
	ResponseBody       []byte
	ResponseStatusCode int
	Debug              bool
}

func New() *Class {
	this := &Class{}
	this.RequestHeader = map[string]string{}
	this.HttpClient = &http.Client{}
	return this
}

func (this *Class) SetProxy(hostport string) error {
	proxyUrl, err := url.Parse("http://" + hostport)
	if err != nil {
		return err
	}
	httpTransport := &http.Transport{}
	this.HttpClient.Transport = httpTransport
	httpTransport.Proxy = http.ProxyURL(proxyUrl)
	return nil
}

func (this *Class) SetSocks5(socks5ipport string) error {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", socks5ipport, nil, proxy.Direct)
	if err != nil {
		//fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		return err
		//os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	this.HttpClient.Transport = httpTransport
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	return nil
}

func (this *Class) NoTransport() {
	this.HttpClient.Transport = nil
}

func (this *Class) RequestByte(method string, sUrl string, query url.Values, body io.Reader, header map[string]string) (int, []byte, error) {
	//以byte数组返回结果
	if query != nil && len(query) > 0 {
		sUrl += "?" + query.Encode()
	}
	this.RequestMethod = method
	this.RequestUrl = sUrl
	this.RequestHeader = header
	this.RequestBody = body
	err := this.Exec()
	return this.ResponseStatusCode, this.ResponseBody, err
}

func (this *Class) SetTimeout(_duration time.Duration) {
	this.HttpClient.Timeout = _duration
}

func (this *Class) Close() {
	this.HttpClient.CloseIdleConnections()
	//this.HttpClient = nil
	//this = nil
}

// client 解析 gzip 返回
func ClientUncompress() {
	client := http.Client{}
	req, err := http.NewRequest("GET", "http://www.baidu.com", nil)
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var buf [1024 * 1024]byte
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return
	}
	_, err = reader.Read(buf[:])
	if err != nil {
		return
	}
	reader.Close()
}

func (this *Class) Exec() error {
	clientReq, err := http.NewRequest(this.RequestMethod, this.RequestUrl, this.RequestBody)
	if err != nil {
		this.ResponseStatusCode = 1
		return err
	}
	if this.RequestHeader != nil {
		for k, v := range this.RequestHeader {
			clientReq.Header.Set(k, v)
		}
	}
	if this.Debug {
		fmt.Println(this.RequestUrl)
		v, _ := json.Marshal(this.RequestHeader)
		fmt.Println(string(v))
	}
	clientRes, err := this.HttpClient.Do(clientReq) //向后端服务器提交数据
	if err != nil {
		this.ResponseStatusCode = 2
		return errors.New("RequestByte-HttpClient-Do:" + err.Error())
	}
	this.ResponseStatusCode = clientRes.StatusCode
	this.ResponseBody, err = ioutil.ReadAll(clientRes.Body) //取得后端服务器返回的数据
	clientRes.Body.Close()
	if err != nil {
		return errors.New("RequestByte-ReadAll:" + err.Error())
	}
	if this.Debug {
		fmt.Println(string(this.ResponseBody))
	}
	return nil
}

func (this *Class) SetRequestJson(v map[string]string) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	this.RequestBody = bytes.NewBuffer(b)
	return nil
}

func (this *Class) SetPost(kv map[string]string) error {
	this.RequestHeader["Content-Type"] = "application/x-www-form-urlencoded"
	u := url.Values{}
	for k, v := range kv {
		u.Add(k, v)
	}
	//fmt.Println(u.Encode())
	this.RequestBody = bytes.NewBuffer([]byte(u.Encode()))
	return nil
}
