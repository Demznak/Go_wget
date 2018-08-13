package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const byteToKByte = 1024

var (
	progress map[string]int64
	mx       sync.Mutex
)

func init() {
	progress = make(map[string]int64)
}

type FileWriter struct {
	Total     int64
	AllLength int64
	Progress  int64
	FileName  string
	Uri       string
}

func (wf *FileWriter) Write(p []byte) (int, error) {
	n := len(p)
	if n > 0 {
		wf.Total += int64(n)
		percentage := float64(wf.Total) / float64(wf.AllLength) * float64(100)
		if int64(percentage)-wf.Progress > 0 {
			wf.Progress = int64(percentage)
			mx.Lock()
			progress[wf.FileName] = wf.Progress
			mx.Unlock()
		}
	}

	return n, nil
}

func DownloadFile(fileParam *FileWriter) {
	response, err := http.Get(fileParam.Uri)
	if err != nil {
		fmt.Printf("\n\t\terr download: %s\n", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("\n\t\tresponse status: %s for url - %s\n\n", response.Status, fileParam.Uri)
		return
	}

	fileParam.AllLength = response.ContentLength

	out, err := os.Create(fileParam.FileName)
	defer out.Close()

	if err != nil {
		fmt.Printf("create file error: %s \n", err.Error())
	}

	fileSize, err := io.Copy(out, io.TeeReader(response.Body, fileParam))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n\t\t\t File %s transferred (%.1f KB).\n\n", fileParam.FileName, float64(fileSize)/byteToKByte)
}

func PrintProgress(bEndDownload chan bool) {
	files := []string{}
	mx.Lock()
	for key, _ := range progress {
		files = append(files, key)
	}
	mx.Unlock()

	fmt.Println(strings.Join(files, " | "))

	for {
		time.Sleep(time.Second)
		strProgress := ""
		for _, val := range files {
			mx.Lock()
			if progressInt, ok := progress[val]; ok {
				strProgress += strconv.FormatInt(progressInt, 10) + "% "
			}
			mx.Unlock()
		}
		fmt.Println(strProgress)
		select {
		case <-bEndDownload:
			return
		default:
			continue
		}
	}
}

func main() {
	var wg sync.WaitGroup

	urlCount := len(os.Args)
	if urlCount <= 1 {
		fmt.Println("not url")
		return
	}
	urlCount -= 1
	urls := os.Args[1:]
	fmt.Printf("You have %d url\n", urlCount)

	wg.Add(urlCount)
	for _, val := range urls {
		fileName := path.Base(val)
		if fileName == "" || fileName == "." {
			fileName = "download_" + time.Now().String()
		}

		mx.Lock()
		progress[fileName] = 0
		mx.Unlock()

		go func(filesParam *FileWriter) {
			defer wg.Done()
			//fmt.Println("Start download url - ", filesParam.Uri)
			DownloadFile(filesParam)
		}(&FileWriter{FileName: fileName, Uri: val})
	}
	bEndDownload := make(chan bool)
	go PrintProgress(bEndDownload)
	wg.Wait()
	bEndDownload <- true
}
