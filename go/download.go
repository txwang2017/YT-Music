package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type AudioMeta struct {
	Url              string
	Id               string
	AudioSampleRate  string
	MimeType         string
	AverageBitrate   int64
	ContentLength    string
	ApproxDurationMs string
}

type QueryRequest struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type YouTubeAudio struct {
	audioMeta *AudioMeta
}

type downloadStatus struct {
	err        error
	readLength int
}

func (youTubeAudio *YouTubeAudio) GetAudioMeta(url string) {
	query := QueryRequest{
		Id:  GetUUID(),
		Url: url,
	}
	request, _ := json.Marshal(query)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(request)
	buff := make([]byte, 1024)
	metaData := make([]byte, 0)
	for n, err := conn.Read(buff); n != 0 && err == nil; n, err = conn.Read(buff) {
		metaData = append(metaData, buff[:n]...)
	}
	audioMeta := AudioMeta{}
	json.Unmarshal(metaData, &audioMeta)
	youTubeAudio.audioMeta = &audioMeta
}

func displayProgress(length float64, status chan downloadStatus, wait *sync.WaitGroup) {
	currLength := 0
	for {
		currStatus := <-status
		n, err := currStatus.readLength, currStatus.err
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("100% downloaded")
			}
			break
		}
		currLength += n
		progress := float64(currLength) / length * 100
		fmt.Printf("%.1f%% downloaded \n", progress)
	}
	wait.Done()
}

func (youTubeAudio *YouTubeAudio) Download(fileName string) {
	resp, _ := http.Get(youTubeAudio.audioMeta.Url)
	buff := make([]byte, 10240)
	filePath := filepath.Join(GetCurrDir(), fileName)
	file, err := os.Create(filePath)
	wait := sync.WaitGroup{}

	if err != nil {
		fmt.Println("Failed to create file")
		return
	}
	status := make(chan downloadStatus)
	length, err := strconv.ParseFloat(youTubeAudio.audioMeta.ContentLength, 64)
	fmt.Println(length, err)
	if err == nil {
		wait.Add(1)
		go displayProgress(length, status, &wait)
	}
	for {
		n, err := resp.Body.Read(buff)
		status <- downloadStatus{err: err, readLength: n}
		if err != nil {
			break
		}
		file.Write(buff[:n])
	}
	file.Close()
	wait.Wait()
}
