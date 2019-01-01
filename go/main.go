package main

import (
	"sync"
)

func main() {
	urls, fileNames, err := ParseCommandLine()
	if err != nil {
		return
	}
	wait := sync.WaitGroup{}
	wait.Add(len(urls))
	for i := 0; i < len(urls); i++ {
		go func(url, fileName string, wait *sync.WaitGroup) {
			y := YouTubeAudio{}
			y.GetAudioMeta(url)
			y.Download(fileName)
			wait.Done()
		}(urls[i], fileNames[i], &wait)
	}
	wait.Wait()
}
