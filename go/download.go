package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type YoutubeAudioMeta struct {
	Url              string
	Id               string
	AudioSampleRate  string
	MimeType         string
	AverageBitrate   int64
	ContentLength    string
	ApproxDurationMs string
	Title            string
}

type QueryRequest struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type downloadStatus struct {
	err        error
	readLength int
}

type DownloadJob struct {
	Urls      []string
	FileNames []string
}

func NewDownloadJob() (*DownloadJob, error) {
	job := DownloadJob{}
	err := job.init()
	if err == nil {
		return &job, nil
	}
	return nil, err
}

func (job *DownloadJob) init() error {
	err := job.parseCommandLine()
	return err
}

func (job *DownloadJob) validateUrl(url string) bool {
	pattern, _ := regexp.Compile(`https://www.youtube.com/watch\?v=[A-Za-z0-9-_]{11}`)
	return pattern.Match([]byte(url))
}

func (job *DownloadJob) parseCommandLine() error {
	var url, fileName, urlListFile string
	urls := make([]string, 0)
	fileNames := make([]string, 0)
	flag.StringVar(&url, "url", "", "URL of youtube video")
	flag.StringVar(&fileName, "name", fmt.Sprintf("%s.mp3", uuid.New().String()), "File name of downloaded audio")
	flag.StringVar(&urlListFile, "list", "", "File path that stored all the youtube video url")
	flag.Parse()
	if url == "" && urlListFile == "" {
		return errors.New("Either url or list file is required")
	}
	if urlListFile != "" {
		buff := make([]byte, 1024)
		data := make([]byte, 0)
		file, err := os.Open(urlListFile)
		defer file.Close()
		if err != nil {
			return err
		}
		for {
			n, err := file.Read(buff)
			if err == nil {
				data = append(data, buff[:n]...)
			} else if err.Error() == "EOF" {
				break
			} else {
				return err
			}
		}
		urlsRaw := strings.Split(string(data), "\n")
		for _, urlRaw := range urlsRaw {
			if job.validateUrl(urlRaw) {
				urls = append(urls, urlRaw)
				fileNames = append(fileNames, "")
			}
		}
	} else {
		urls = append(urls, url)
		fileNames = append(fileNames, fileName)
	}
	job.Urls = urls
	job.FileNames = fileNames
	return nil
}

func (job *DownloadJob) getAudioMeta(url string) (*YoutubeAudioMeta, error) {
	query := QueryRequest{
		Id:  GetUUID(),
		Url: url,
	}
	request, _ := json.Marshal(query)
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		return nil, err
	}
	conn.Write(request)
	buff := make([]byte, 1024)
	metaData := make([]byte, 0)
	for n, err := conn.Read(buff); n != 0 && err == nil; n, err = conn.Read(buff) {
		metaData = append(metaData, buff[:n]...)
	}
	audioMeta := YoutubeAudioMeta{}
	json.Unmarshal(metaData, &audioMeta)
	return &audioMeta, nil
}

func displayProgress(length float64, status chan downloadStatus, sequence int, wait *sync.WaitGroup) {
	currLength := 0
	for {
		currStatus := <-status
		n, err := currStatus.readLength, currStatus.err
		fmt.Printf("\033[%d;0H", sequence)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("\r 100% downloaded		")
			}
			break
		}
		currLength += n
		progress := float64(currLength) / length * 100
		fmt.Printf("\r%.1f%% downloaded		", progress)
	}
	wait.Done()
}

func (job *DownloadJob) download(fileName string, audioMeta *YoutubeAudioMeta, sequence int) error {
	resp, _ := http.Get(audioMeta.Url)
	buff := make([]byte, 10240)
	if fileName == "" {
		fileName = audioMeta.Title
	}
	filePath := filepath.Join(GetMusicDir(), fileName)
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	wait := sync.WaitGroup{}
	status := make(chan downloadStatus)
	length, err := strconv.ParseFloat(audioMeta.ContentLength, 64)
	if err == nil {
		wait.Add(1)
		go displayProgress(length, status, sequence, &wait)
	}
	for {
		n, err := resp.Body.Read(buff)
		status <- downloadStatus{err: err, readLength: n}
		if err != nil {
			fmt.Println(err)
			break
		}
		file.Write(buff[:n])
	}
	wait.Wait()
	return nil
}

func (job *DownloadJob) Download() {
	wait := sync.WaitGroup{}
	wait.Add(len(job.Urls))
	for i := 0; i < len(job.Urls); i++ {
		go func(url, fileName string, sequence int, wait *sync.WaitGroup) {
			audioMeta, err := job.getAudioMeta(url)
			if err != nil {
				fmt.Println(err)
				wait.Done()
				return
			}
			err = job.download(fileName, audioMeta, sequence)
			if err != nil {
				fmt.Println(err)
				wait.Done()
				return
			}
			wait.Done()
		}(job.Urls[i], job.FileNames[i], i, &wait)
	}
	wait.Wait()
}
