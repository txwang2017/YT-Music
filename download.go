package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func Download(link *VideoLink) {
	resp, _ := http.Get(link.Url)
	buff := make([]byte, 10240)
	filePath := filepath.Join(GetCurrDir(), "download.mp3")
	file, _ := os.Create(filePath)
	for n, _ := resp.Body.Read(buff); n != 0; n, _ = resp.Body.Read(buff) {
		file.Write(buff[:n])
	}
	file.Close()
}
