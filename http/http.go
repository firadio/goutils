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

type HttpStruct struct {
	HttpClient *http.Client
}

func New() *HttpStruct {
	http1 := &HttpStruct{}
	http1.HttpClient = &http.Client{}
	return http1
}

func (http1 *HttpStruct) SetSocks5(socks5ipport string) error {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", socks5ipport, nil, proxy.Direct)
	if err != nil {
		//fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		return err
		//os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	http1.HttpClient.Transport = httpTransport
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	return nil
}

func (http1 *HttpStruct) NoTransport() {
	http1.HttpClient.Transport = nil
}

func (http1 *HttpStruct) RequestByte(method string, sUrl string, query url.Values, body io.Reader, header map[string]string) (int, []byte, error) {
	//以byte数组返回结果
	if query != nil {
		sUrl += "?" + query.Encode()
	}
	clientReq, err := http.NewRequest(method, sUrl, body)
	if err != nil {
		return 0, nil, err
	}
	if header != nil {
		for k, v := range header {
			clientReq.Header.Set(k, v)
		}
	}
	clientRes, err := http1.HttpClient.Do(clientReq) //向后端服务器提交数据
	if err != nil {
		return 0, nil, errors.New("RequestByte-HttpClient-Do:" + err.Error())
	}
	clientResBody, err := ioutil.ReadAll(clientRes.Body) //取得后端服务器返回的数据
	clientRes.Body.Close()
	if err != nil {
		return clientRes.StatusCode, nil, errors.New("RequestByte-ReadAll:" + err.Error())
	}
	return clientRes.StatusCode, clientResBody, nil
}

func (http1 *HttpStruct) SetTimeout(_duration time.Duration) {
	http1.HttpClient.Timeout = _duration
}

func (http1 *HttpStruct) Close() {
	http1.HttpClient.CloseIdleConnections()
	http1.HttpClient = nil
	http1 = nil
}
