package utils

import (
	"encoding/base64"
	"errors"
	"unsafe"

	"github.com/xxtea/xxtea-go/xxtea"
)

func XxteaEncrypt(clientResBodyNew []byte, key string) string {
	clientBodyEncrypt := xxtea.Encrypt(clientResBodyNew, []byte(key))
	clientBodyBase64 := base64.StdEncoding.EncodeToString(clientBodyEncrypt)
	return clientBodyBase64
}

func XxteaBase64ToByte(base64text string, key string) ([]byte, error) {
	if base64text == "" {
		//base64text是空的
		return nil, errors.New("Base64TextIsEmpty")
	}
	bEncrypted, err := base64.StdEncoding.DecodeString(base64text)
	if err != nil {
		//base64解码失败
		return nil, errors.New("base64DecodeError")
	}
	bDecrypted := xxtea.Decrypt(bEncrypted, []byte(key))
	if len(bDecrypted) == 0 {
		//xxtea解密失败
		return nil, errors.New("DecryptError")
	}
	return bDecrypted, nil
}

func XxteaDecrypt(base64text string, key string) (string, error) {
	bDecrypted, err := XxteaBase64ToByte(base64text, key)
	if err != nil {
		return "", err
	}
	str := (*string)(unsafe.Pointer(&bDecrypted))
	return *str, nil
}
