package main

// "golang.org/x/net/html"

func main() {
	id := "VC2rAxRID9s"
	youtube := NewYouTubeRequest(id)
	resp := youtube.GetResponse()
	dom := GetYouTubeDOM(resp)
	dom.Parse(dom.root, "div", "id", "player-wrap", 0)
	dom.Parse(dom.currNode, "script", "", "", 3)
	links := dom.GetLinks()
	Download(links[1])
}
