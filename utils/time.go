package utils

import (
	"time"
)

func TimeNow() time.Time {
	cstZone := time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone)
}

func TimestampFloat64() float64 {
	// 返回是秒级时间戳（带小数点的）
	return float64(TimeNow().UnixNano()) / 1e9
}

func TimestampInt32() int {
	// 返回是秒级时间戳（int32型，不带小数点的，PHP标准）
	return int(TimeNow().UnixNano() / 1e9)
}

func TimestampInt64() int64 {
	// 返回是秒级时间戳（long型，不带小数点的）
	return TimeNow().UnixNano() / 1e9
}

func TimestampMsInt64() int64 {
	// 返回是毫秒级时间戳（long型，不带小数点的，JAVA/JavaScript标准）
	return TimeNow().UnixNano() / 1e6
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
