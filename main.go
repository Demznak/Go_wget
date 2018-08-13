package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func DownloadFile(fileParam *FileParams) {
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
	fileParam.Reader = response.Body

	out, err := os.Create(fileParam.FileName)
	defer out.Close()

	if err != nil {
		fmt.Printf("create file error: %s \n", err.Error())
	}

	time.Sleep(time.Second * 2)
	size, err := io.Copy(out, fileParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n\t\t\t%s Transferred. (%.1f KB)\n\n", fileParam.FileName, float64(size)/byteToKbyte)
}

func PrintProgress(bEndDownload chan bool) {
	files := mapFilesProgress.GetFileNames()
	fmt.Println(strings.Join(files, " | "))

	for {
		time.Sleep(time.Second)
		strProgress := ""
		for _, val := range files {
			progressInt, ok := mapFilesProgress.GetProgressFile(val)
			if !ok {
				progressInt = 0
			}
			//strProgress += string(progress) + "% "
			strProgress += strconv.FormatInt(progressInt, 10) + "% "
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
		fileName := getFilename(val)
		mapFilesProgress.SetProgressFile(fileName, 0)

		go func(filesParam *FileParams) {
			defer wg.Done()
			//fmt.Println("Start download url - ", filesParam.Uri)
			DownloadFile(filesParam)
		}(&FileParams{FileName: fileName, Uri: val})
	}
	bEndDownload := make(chan bool)
	go PrintProgress(bEndDownload)
	wg.Wait()
	bEndDownload <- true
}
