package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func HttpRequestByte(method string, sUrl string, query url.Values, body io.Reader, httpHeader http.Header) (int, []byte, error) {
	//以byte数组返回结果
	if query != nil {
		sUrl += "?" + query.Encode()
	}
	clientReq, err := http.NewRequest(method, sUrl, body)
	if err != nil {
		return 0, nil, err
	}
	if httpHeader != nil {
		clientReq.Header = httpHeader
	} else {
		clientReq.Header.Set("user-agent", "")
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
