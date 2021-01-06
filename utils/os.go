package utils

import (
	"os"
)

func IsExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		//没报错，说明文件肯定存在
		return true
	}
	//有报错，说明文件可能不存在
	if os.IsNotExist(err) {
		//文件确实不存在
		return false
	}
	//文件可能存在，或有其他问题
	return true
}
