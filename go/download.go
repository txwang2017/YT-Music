package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// {"lastModified": "1540116568839208",
// "quality": "tiny", "bitrate": 128319,
//   "id": 19, "itag": 140, "audioQuality": "AUDIO_QUALITY_MEDIUM",
//    "contentLength": "3782204", "mimeType": "audio/mp4; codecs=\"mp4a.40.2\"",
// 	"approxDurationMs": "238097", "averageBitrate": 127081, "highReplication": true,
// 	"initRange": {"start": "0", "end": "591"},
// 	"url": "https://r3---sn-n4v7sn7l.googlevideo.com/videoplayback?dur=238.097&c=WEB&fvip=3&lmt=1540116568839208&clen=3782204&id=o-ABfV2Y1y5LR-mhj6mTMS4l78v9ab5ZSTvp6-tPDj2P5w&mn=sn-n4v7sn7l%2Csn-n4v7knls&mm=31%2C29&ip=2601%3A647%3A4000%3Ad7f7%3Aed4c%3A2ec7%3Af0b6%3A243b&ms=au%2Crdu&ei=Pd8qXPLvLMmNkgaNupmYDw&pl=34&mv=m&ipbits=0&sparams=clen%2Cdur%2Cei%2Cgir%2Cid%2Cinitcwndbps%2Cip%2Cipbits%2Citag%2Ckeepalive%2Clmt%2Cmime%2Cmm%2Cmn%2Cms%2Cmv%2Cpl%2Crequiressl%2Csource%2Cexpire&signature=2BC1B692618F01D33053BFE069340EBAD5564767.D65065E30962458DFC6CDBAAF9314632895C8041&source=youtube&txp=5432432&itag=140&expire=1546335133&key=yt6&mime=audio%2Fmp4&gir=yes&keepalive=yes&requiressl=yes&mt=1546313415&initcwndbps=1930000",
// "audioSampleRate": "44100", "projectionType": "RECTANGULAR"}

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
