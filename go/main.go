package main

import (
	"fmt"
)

func main() {
	downloadJob, err := NewDownloadJob()
	if err != nil {
		fmt.Println(err)
		return
	}
	downloadJob.Download()
}
