package main

import (
	"net/http"
	"os"
)

func Download(link *VideoLink) {
	resp, _ := http.Get(link.Url)
	buff := make([]byte, 10240)
	file, _ := os.Create("/Users/Tianxiong.wang/test/gozoo/download.mp3")
	for n, _ := resp.Body.Read(buff); n != 0; n, _ = resp.Body.Read(buff) {
		file.Write(buff[:n])
	}
	file.Close()
}
