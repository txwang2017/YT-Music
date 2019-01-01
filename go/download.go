package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
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
	fmt.Println(string(metaData))
	json.Unmarshal(metaData, &audioMeta)
	fmt.Println(audioMeta)
	youTubeAudio.audioMeta = &audioMeta
}

func (youTubeAudio *YouTubeAudio) Download(fileName string) {
	resp, _ := http.Get(youTubeAudio.audioMeta.Url)
	buff := make([]byte, 10240)
	filePath := filepath.Join(GetCurrDir(), fileName)
	file, _ := os.Create(filePath)
	for n, err := resp.Body.Read(buff); n != 0 && err == nil; n, err = resp.Body.Read(buff) {
		file.Write(buff[:n])
	}
	file.Close()
}
