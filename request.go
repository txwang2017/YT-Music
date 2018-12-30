package main

import (
	"fmt"
	"io"
	"net/http"
)

//YouTube is a wrapper of youtube client
type YouTube struct {
	client  *http.Client
	request *http.Request
}

func (youtube *YouTube) sendRequest() *http.Response {
	resp, err := youtube.client.Do(youtube.request)
	if err != nil {
		return nil
	}
	return resp
}

//ReadResponse reads the resp as a string
func (youtube *YouTube) ReadResponse() string {
	resp := youtube.sendRequest()
	if resp == nil {
		return ""
	}
	buff := make([]byte, 1024)
	data := make([]byte, 0)
	for n, _ := resp.Body.Read(buff); n != 0; n, _ = resp.Body.Read(buff) {
		data = append(data, buff...)
	}
	res := string(data[:])
	return res
}

//GetResponse return the body of response
func (youtube *YouTube) GetResponse() io.Reader {
	resp := youtube.sendRequest()
	if resp == nil {
		return nil
	}
	return resp.Body
}

//NewYouTubeRequest creates a youtube client
func NewYouTubeRequest(videoID string) *YouTube {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// request.Header.Set("accept-encoding", "gzip, deflate, br")
	request.Header.Set("accept-language", "en-US,en;q=0.9")
	request.Header.Set("cookie", "VISITOR_INFO1_LIVE=7gcjyw3ehOI; CONSENT=YES+US.en+20170326-06-0; PREF=f1=50000000; GPS=1; YSC=lTe68erqg-M")
	request.Header.Set("upgrade-insecure-requests", "1")
	request.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	youtube := YouTube{
		client:  client,
		request: request,
	}
	return &youtube
}
