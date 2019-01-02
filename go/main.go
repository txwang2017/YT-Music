package main

import (
	"sync"
)

func main() {
	downloadJob, err := NewDownloadJob()
	if err != nil {
		return
	}
	wait := sync.WaitGroup{}
	wait.Add(len(downloadJob.Urls))
	for i := 0; i < len(downloadJob.Urls); i++ {
		go func(url, fileName string, wait *sync.WaitGroup) {
			y := YouTubeAudio{}
			y.GetAudioMeta(url)
			y.Download(fileName)
			wait.Done()
		}(downloadJob.Urls[i], downloadJob.FileNames[i], &wait)
	}
	wait.Wait()
}
