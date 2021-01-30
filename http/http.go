package http

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type Class struct {
	HttpClient    *http.Client
	RequestUrl    string
	RequestHeader map[string]string
	ResponseBody  []byte
}

func New() *Class {
	this := &Class{}
	this.HttpClient = &http.Client{}
	return this
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
	this.RequestUrl = sUrl
	this.RequestHeader = header
	clientReq, err := http.NewRequest(method, sUrl, body)
	if err != nil {
		return 0, nil, err
	}
	if header != nil {
		for k, v := range header {
			clientReq.Header.Set(k, v)
		}
	}
	clientRes, err := this.HttpClient.Do(clientReq) //向后端服务器提交数据
	if err != nil {
		return 0, nil, errors.New("RequestByte-HttpClient-Do:" + err.Error())
	}
	clientResBody, err := ioutil.ReadAll(clientRes.Body) //取得后端服务器返回的数据
	this.ResponseBody = clientResBody
	clientRes.Body.Close()
	if err != nil {
		return clientRes.StatusCode, nil, errors.New("RequestByte-ReadAll:" + err.Error())
	}
	return clientRes.StatusCode, clientResBody, nil
}

func (this *Class) SetTimeout(_duration time.Duration) {
	this.HttpClient.Timeout = _duration
}

func (this *Class) Close() {
	this.HttpClient.CloseIdleConnections()
	this.HttpClient = nil
	this = nil
}
