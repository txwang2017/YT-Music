package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Need video URL and target file name")
		return
	}
	videoLink := os.Args[1]
	fileName := os.Args[2]
	youtube := NewYouTubeRequest(videoLink)
	resp := youtube.GetResponse()
	dom := GetYouTubeDOM(resp)
	dom.Parse(dom.root, "div", "id", "player-wrap", 0)
	dom.Parse(dom.currNode, "script", "", "", 3)
	links := dom.GetLinks()
	Download(links[0], fileName)
}
