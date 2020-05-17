package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	var files []string
	rootDir := flag.String("rootDir", "", "the path of directory contains .vtt files")
	flag.Parse()

	if *rootDir == "" {
		fmt.Print(" Please input VTT directory path!")
		return
	}

	//Iterate the specify directory  and find all .vtt files
	err := filepath.Walk(*rootDir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".vtt") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		converterSubtitles(file)
	}

}

//converter the vtt subtitle to the srt subtitle
func converterSubtitles(path string) {
	//获取文件所在的文件夹
	dir, fileName := filepath.Split(path)

	srtFile := dir + strings.TrimSuffix(fileName, ".vtt") + ".srt"
	f, err := os.Create(srtFile)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	defer f.Close()

	vttFile, err := os.Open(path)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	defer vttFile.Close()

	r := bufio.NewReader(vttFile)
	lineNumber := 1
	var lineStr string
	var isCaptionLine = false
	for {
		//read file by line
		buf, _, err := r.ReadLine()

		if err != nil {
			if err == io.EOF {
				break
			}
		}

		lineStr = string(buf)
		if isTimeLine(lineStr) {

			//字幕行号
			f.WriteString(strconv.Itoa(lineNumber) + "\n")
			lineNumber++


			timeStrArr := strings.Split(strings.TrimSuffix(lineStr, "\n"), " --> ")
			const (
				vttFormatStr = "04:05.000"
				srtFormatStr = "00:04:05.000"
			)
			time1, _ := time.Parse(vttFormatStr, timeStrArr[0])
			timeSrt1 := strings.Replace(time1.Format(srtFormatStr), ".", ",", 1)
			time2, _ := time.Parse(vttFormatStr, timeStrArr[1])
			timeSrt2 := strings.Replace(time2.Format(srtFormatStr), ".", ",", 1)
			lineStr = timeSrt1 + " --> " + timeSrt2 + "\n"
			f.WriteString(lineStr)
			//下一行就是字幕
			isCaptionLine = true

		} else if isCaptionLine {
			f.WriteString(lineStr + "\n\n")
			isCaptionLine = false
		}
	}
}

/**
  判断是否是时间行
*/
func isTimeLine(line string) bool {
	return strings.Contains(line, "-->")
}
