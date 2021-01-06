package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

var pwdpre_key = "Vn3gAlid0w3"

type HashrateQuestion struct {
	Pwdpre     string `json:"pwdpre"`     //密码前缀
	Crypto     string `json:"crypto"`     //算法
	RangeStart uint64 `json:"rangeStart"` //范围开始
	RangeEnd   uint64 `json:"rangeEnd"`   //范围结束
	CheckCode  string `json:"checkCode"`  //校验码
}

type HashrateAnswer struct {
	Pwdpre     string `json:"pwdpre"`     //密码前缀
	Crypto     string `json:"crypto"`     //算法
	RangeStart uint64 `json:"rangeStart"` //范围开始
	RangeEnd   uint64 `json:"rangeEnd"`   //范围结束
	Subpwd     uint64 `json:"subpwd"`     //子密码
}

func getHashrateAnswerFromInterface(interfaceHashrate interface{}) (HashrateAnswer, error) {
	oHashrateAnswer := HashrateAnswer{}
	mHashrate, ok := interfaceHashrate.(map[string]interface{})
	if !ok {
		return oHashrateAnswer, errors.New("提供的common_hashrate无法转成map[string]interface{}")
	}
	//pwdpre
	ifPwdpre, ok := mHashrate["pwdpre"]
	if !ok {
		return oHashrateAnswer, errors.New("pwdpre不存在")
	}
	oHashrateAnswer.Pwdpre, ok = ifPwdpre.(string)
	if !ok {
		return oHashrateAnswer, errors.New("pwdpre格式不对")
	}
	//rangeStart
	ifRangeStart, ok := mHashrate["rangeStart"]
	if !ok {
		return oHashrateAnswer, errors.New("rangeStart不存在")
	}
	rangeStart, ok := ifRangeStart.(float64)
	if !ok {
		return oHashrateAnswer, errors.New("rangeStart格式不对")
	}
	oHashrateAnswer.RangeStart = uint64(rangeStart)
	//rangeEnd
	ifRangeEnd, ok := mHashrate["rangeEnd"]
	if !ok {
		return oHashrateAnswer, errors.New("rangeEnd不存在")
	}
	rangeEnd, ok := ifRangeEnd.(float64)
	if !ok {
		return oHashrateAnswer, errors.New("rangeEnd格式不对")
	}
	oHashrateAnswer.RangeEnd = uint64(rangeEnd)
	//subpwd
	ifSubpwd, ok := mHashrate["subpwd"]
	if !ok {
		return oHashrateAnswer, errors.New("subpwd不存在")
	}
	subpwd, ok := ifSubpwd.(float64)
	if !ok {
		return oHashrateAnswer, errors.New("subpwd格式不对")
	}
	oHashrateAnswer.Subpwd = uint64(subpwd)
	return oHashrateAnswer, nil
}

func getHashrateAnswerFromString(sHashrate string) (HashrateAnswer, error) {
	oHashrateAnswer := HashrateAnswer{}
	err := json.Unmarshal([]byte(sHashrate), &oHashrateAnswer)
	return oHashrateAnswer, err
}

func HashrateCheck(userReqJson map[string]interface{}) error {
	interfaceHashrate, ok := userReqJson["common_hashrate"]
	if !ok {
		return errors.New("没有提供common_hashrate")
	}
	sHashrate, ok := interfaceHashrate.(string)
	if !ok {
		return errors.New("提供的common_hashrate无法转成string")
	}
	oHashrateAnswer, err := getHashrateAnswerFromString(sHashrate)
	if err != nil {
		return err
	}
	//参数全部获得，开始比对密钥
	uint64Subpwd := getSubPwd(oHashrateAnswer.Pwdpre, oHashrateAnswer.RangeStart, oHashrateAnswer.RangeEnd)
	if uint64Subpwd != uint64(oHashrateAnswer.Subpwd) {
		return errors.New("提供的答案不正确")
	}
	return nil
}

func HashrateGetNew(crypto string) HashrateQuestion {
	upperCrypto := strings.ToUpper(crypto)
	pwdpreInfo := HashrateQuestion{}
	pwdpreInfo.Crypto = upperCrypto
	b := make([]byte, 16)
	rand.Read(b)
	pwdpreInfo.Pwdpre = ByteToCRC64DecStr(b)
	pwdpreInfo.RangeStart = 10000
	pwdpreInfo.RangeEnd = 99999
	uint64Subpwd := getSubPwd(pwdpreInfo.Pwdpre, pwdpreInfo.RangeStart, pwdpreInfo.RangeEnd)
	stringSubpwd := strconv.FormatUint(uint64Subpwd, 10)
	checkCodeText := pwdpreInfo.Pwdpre + stringSubpwd
	switch upperCrypto {
	case "CRC32":
		pwdpreInfo.CheckCode = CRC32(checkCodeText)
	case "CRC64":
		pwdpreInfo.CheckCode = StrToCRC64DecStr(checkCodeText)
	case "MD5":
		pwdpreInfo.CheckCode = StrToMD5HexStr(checkCodeText)
	case "SHA1":
		pwdpreInfo.CheckCode = SHA1(checkCodeText)
	}
	return pwdpreInfo
}

func getSubPwd(Pwdpre string, RangeStart uint64, RangeEnd uint64) uint64 {
	str := Pwdpre + pwdpre_key
	crc64Ret := StrToCRC64Uint64(str)
	div := RangeEnd - RangeStart
	if div == 0 {
		return RangeStart
	}
	if div < 0 {
		return (crc64Ret % (-div + 1)) + RangeEnd
	}
	return (crc64Ret % (div + 1)) + RangeStart
}
