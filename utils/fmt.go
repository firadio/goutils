package utils

import (
	"fmt"
)

func FmtPrintf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func FmtPrintln(args ...interface{}) {
	now := TimeNow()
	datetime := now.Format("2006-01-02 15:04:05")
	args1 := make([]interface{}, 0)
	args1 = append(args1, datetime)
	for _, v1 := range args {
		args1 = append(args1, v1)
	}
	fmt.Println(args1...)
}
