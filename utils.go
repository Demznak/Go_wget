package main

import (
	"fmt"
	"gowget/progress"
	"io"
	"net/url"
	"path"
	"strings"
	"time"
)

const byteToKbyte = 1024.0

var mapFilesProgress *progress.Progress

func init() {
	mapFilesProgress = progress.GetProgress()
}

type FileParams struct {
	io.Reader
	Total     int64
	AllLength int64
	Progress  int64
	FileName  string
	Uri       string
}

func getFilename(url string) string {
	fileName := path.Base(url)

	if fileName == "" || fileName == "." {
		return "download_" + time.Now().String()
	}

	return fileName
}

func cutBefore(s, sep string) string {
	if strings.Contains(s, sep) {
		return strings.Split(s, sep)[1]
	}

	return s
}

func checkUri(uri string) bool {
	_, e := url.ParseRequestURI(uri)
	if e != nil {
		fmt.Println("invalid url")
		return false
	}
	return true
}

func (fp *FileParams) Read(p []byte) (int, error) {
	n, err := fp.Reader.Read(p)
	if n > 0 {
		fp.Total += int64(n)
		percentage := float64(fp.Total) / float64(fp.AllLength) * float64(100)

		if int64(percentage)-fp.Progress > 0 {
			//fmt.Fprintf(os.Stderr, is)
			fp.Progress = int64(percentage)
			mapFilesProgress.SetProgressFile(fp.FileName, fp.Progress)
		}
	}

	return n, err
}
