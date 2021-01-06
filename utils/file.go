package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type FileReadCB func([]string)

func FileReadAllLine(filename string, fun FileReadCB) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		sLine, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		sLine = strings.TrimRight(sLine, "\r\n")
		aLine := strings.Split(sLine, "----")
		fun(aLine)
	}
}
