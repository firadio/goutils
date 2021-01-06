package utils

import (
	"encoding/json"
	"errors"
	"strings"
	"unsafe"
)

type ResponseJson struct {
	Time     float64                `json:"time"`
	Ret      int                    `json:"ret"`
	Code     string                 `json:"code"`
	Msg      string                 `json:"msg"`
	Msgbox   map[string]interface{} `json:"msgbox"`
	Data     map[string]interface{} `json:"data"`
	Cache    map[string]interface{} `json:"cache"`
	Setting  map[string]string      `json:"setting"`
	Debug    map[string]interface{} `json:"debug"`
	Hashrate HashrateQuestion       `json:"hashrate"`
}

func GetResJsonNew() ResponseJson {
	apiResJson := ResponseJson{}
	apiResJson.Time = TimestampFloat64()
	apiResJson.Data = map[string]interface{}{}
	return apiResJson
}

func GetDataByJsonRaw(apiResJson ResponseJson) []byte {
	// 下面将 apiResJson 转回字符串并加密传输
	clientResBodyNew, err := json.Marshal(&apiResJson)
	if err != nil {
		// panic(err)
		return nil
	}
	return clientResBodyNew
}

func GetJsonByRes(apiResJson ResponseJson) []byte {
	// 下面将 apiResJson 转回字符串并加密传输
	clientResBodyNew, err := json.Marshal(&apiResJson)
	if err != nil {
		return nil
	}
	return clientResBodyNew
}

type TokenInfo struct {
	Logined   float64 `json:"logined"`
	Actived   float64 `json:"actived"`
	LoginTime string
	UserId    string `json:"user_id"`
	LoginIp   string
	Plat      string
}

func GetCommonToken(userReqJson map[string]interface{}, service string) (TokenInfo, error) {
	ifCommonToken, ok := userReqJson["common_token"]
	oTokenInfo := TokenInfo{}
	if !ok {
		//common_token没提供
		return oTokenInfo, errors.New("common_token is not found")
	}
	common_token := ifCommonToken.(string)
	if common_token == "" {
		//common_token是空串
		return oTokenInfo, errors.New("common_token is empty string")
	}
	key := "qaz@!#321"
	aService := strings.Split(service, ".")
	if aService[0] == "family" {
		key = "DafJ8up2i5I31"
	}
	byteUserTokenDecrypted, err := XxteaBase64ToByte(common_token, key)
	if err != nil {
		//common_token解密失败
		return oTokenInfo, err
	}
	if IsJsonMap(byteUserTokenDecrypted) {
		err := json.Unmarshal(byteUserTokenDecrypted, &oTokenInfo)
		if err != nil {
			return oTokenInfo, err
		}
		return oTokenInfo, nil
	} else {
		sUserTokenDecrypted := (*string)(unsafe.Pointer(&byteUserTokenDecrypted))
		aUserToken := strings.SplitN(*sUserTokenDecrypted, "|", 4)
		iLen := len(aUserToken)
		oTokenInfo.LoginTime = aUserToken[0]
		if iLen > 1 {
			oTokenInfo.UserId = aUserToken[1]
		}
		if iLen > 2 {
			oTokenInfo.LoginIp = aUserToken[2]
		}
		if iLen > 3 {
			oTokenInfo.Plat = aUserToken[3]
		}
	}
	return oTokenInfo, nil
}

func IsJsonMap(byte []byte) bool {
	if len(byte) == 0 {
		return false
	}
	if byte[0] == 1 {
		return true
	}
	return false
}
