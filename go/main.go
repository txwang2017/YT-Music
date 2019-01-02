package main

import "fmt"

func main() {
	downloadJob, err := NewDownloadJob()
	if err != nil {
		return
	}
	downloadJob.Download()
	fmt.Println("a")
}
