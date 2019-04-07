package iplocate

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

//ProcessLine 处理文件中的一行数据, interrupt=true 表示不需要继续读取将停止读取文件，否则继续读取
type ProcessLine func(line string) (interrupt bool)

//FileLoader 文件读取
type FileLoader struct {
}

//ParseFile 读取文件
func (l *FileLoader) ParseFile(filePath string, process ProcessLine) error {
	rw, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("open file :[%v] error:[%v]", filePath, err)
	}
	defer rw.Close()

	rb := bufio.NewReader(rw)
	for {
		line, _, err := rb.ReadLine()
		if err == io.EOF {
			break
		}
		if nil != process {
			if interrupt := process(string(line)); interrupt {
				break
			}
		} else {
			// do nothing
		}
	}
	return nil
}

//LoadFile 加载文件内容
func (l *FileLoader) LoadFile(filePath string) ([]string, error) {
	var lines = make([]string, 0, 100000)

	var err = l.ParseFile(filePath,
		func(line string) bool {
			lines = append(lines, line)
			return false
		})

	return lines, err
}
