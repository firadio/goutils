package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash/crc32"
	"hash/crc64"
	"strconv"
)

func Uint64ToDecStr(ui64 uint64) string {
	return strconv.FormatUint(ui64, 10)
}

var crc32Table = crc32.MakeTable(0xD5828281)

func CRC32(str string) string {
	crc32Int := crc32.Checksum([]byte(str), crc32Table)
	return strconv.FormatUint(uint64(crc32Int), 10)
}

var crc64Table = crc64.MakeTable(0xC96C5795D7870F42)

func ByteToCRC64Uint64(b []byte) uint64 {
	//二进制转CRC64整型（最原始转换）
	return crc64.Checksum(b, crc64Table)
}
func ByteToCRC64DecStr(b []byte) string {
	//二进制转CRC64以十进制字符串表示（获得随机密码用）
	return Uint64ToDecStr(crc64.Checksum(b, crc64Table))
}
func StrToCRC64Uint64(str string) uint64 {
	//字符串转CRC64整型（通过前缀获得一个容易计算的位）
	return ByteToCRC64Uint64([]byte(str))
}
func StrToCRC64DecStr(str string) string {
	//字符串转CRC64以十进制字符串表示（最常用）
	return Uint64ToDecStr(StrToCRC64Uint64(str))
}

func ByteToMD5Byte(bytes []byte) [16]byte {
	return md5.Sum(bytes)
}
func ByteToMD5HexStr(bytes []byte) string {
	Bytes := md5.Sum(bytes)
	return hex.EncodeToString(Bytes[:])
}
func StrToMD5HexStr(str string) string {
	Bytes := md5.Sum([]byte(str))
	return hex.EncodeToString(Bytes[:])
}

func SHA1(str string) string {
	Bytes := sha1.Sum([]byte(str))
	return hex.EncodeToString(Bytes[:])
}
