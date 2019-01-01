package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func ParseCommandLine() ([]string, []string, error) {
	var url, fileName, urlListFile string
	urls := make([]string, 0)
	fileNames := make([]string, 0)
	flag.StringVar(&url, "url", "", "URL of youtube video")
	flag.StringVar(&fileName, "name", fmt.Sprintf("%s.mp3", uuid.New().String()), "File name of downloaded audio")
	flag.StringVar(&urlListFile, "list", "", "File path that stored all the youtube video url")
	flag.Parse()
	if url == "" && urlListFile == "" {
		return nil, nil, errors.New("Either url or list file is required")
	}
	if urlListFile != "" {
		buff := make([]byte, 1024)
		data := make([]byte, 0)
		file, err := os.Open(urlListFile)
		if err != nil {
			fmt.Println(err, "*****")
			return nil, nil, err
		}
		for {
			n, err := file.Read(buff)
			if err == nil {
				data = append(data, buff[:n]...)
			} else if err.Error() == "EOF" {
				break
			} else {
				return nil, nil, err
			}
		}
		urls = strings.Split(string(data), "\n")
		for i := 0; i < len(urls); i++ {
			fileNames = append(fileNames, fmt.Sprintf("%s.mp3", uuid.New().String()))
		}
	} else {
		urls = append(urls, url)
		fileNames = append(fileNames, fileName)
	}
	return urls, fileNames, nil
}

//CompareBytes check if the target string is equal to string represents by byte array
func CompareBytes(source []byte, target string) bool {
	s := string(source[:])
	if target == s {
		return true
	}
	return false
}

//CompareSlice compare if two slices are equal to each other
func CompareSlice(source []byte, target []byte) bool {
	if len(source) != len(target) {
		return false
	}
	for i := 0; i < len(source); i++ {
		if source[i] != target[i] {
			return false
		}
	}
	return true
}

func GetCurrDir() string {
	path := os.Args[0]
	path, _ = filepath.Abs(path)
	currPath := filepath.Dir(path)
	return currPath
}

func GetUUID() string {
	uid := uuid.New()
	return uid.String()
}

func GetMusicDir() string {
	path := os.Getenv("HOME")
	path = filepath.Join(path, "Music")
	return path
}
