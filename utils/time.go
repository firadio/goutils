package utils

import (
	"time"
)

func TimeNow() time.Time {
	cstZone := time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone)
}

func TimestampFloat64() float64 {
	return float64(TimeNow().UnixNano()) / 1e9
}

func TimestampInt64() int64 {
	return TimeNow().UnixNano() / 1e9
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
