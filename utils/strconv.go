package utils

import (
	"strconv"
)

func InterfaceToString(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ""
}

func itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}
