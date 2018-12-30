package main

import (
	"encoding/json"
	"io"

	"golang.org/x/net/html"
)

type KVPair struct {
	key []byte
	val []byte
}

func newKVPair(key, val string) *KVPair {
	return &KVPair{key: []byte(key), val: []byte(val)}
}

var URLFormater [](*KVPair) = [](*KVPair){
	newKVPair(`\\\"`, `'`),
	newKVPair(`\\u0026`, `&`),
	newKVPair(`\u0026`, `&`),
	newKVPair(`\"`, `"`),
	newKVPair(`\/`, `/`),
}

//YouTubeDOM is DOM tree of youtube
type YouTubeDOM struct {
	root     *html.Node
	currNode *html.Node
}

type writer struct {
	data *[]byte
}

type VideoLink struct {
	Url          string
	Quality      string
	QualityLabel string
}

func (videoLink *VideoLink) validate() bool {
	res := (videoLink.Url != "" && videoLink.Quality != "" && videoLink.QualityLabel != "")
	return res
}

func (w writer) Write(p []byte) (n int, err error) {
	*w.data = append(*w.data, p...)
	return len(p), nil
}

func dfs(root *html.Node, tagName, attrKey, attrVal string, targetNum int) *html.Node {
	currNum := 1
	for ; root != nil; root = root.NextSibling {
		res := dfs(root.FirstChild, tagName, attrKey, attrVal, targetNum)
		if res != nil {
			return res
		}
		if root.Type == html.ElementNode && root.Data == tagName {
			currNum++
			if attrKey == "" && attrVal == "" && currNum == targetNum {
				return root
			}
			for _, attr := range root.Attr {
				if attr.Key == attrKey && attr.Val == attrVal {
					return root
				}
			}
		}
	}
	return nil
}

//GetYouTubeDOM returns the root of youtube DOM tree
func GetYouTubeDOM(resp io.Reader) *YouTubeDOM {
	root, err := html.Parse(resp)
	if err != nil {
		return nil
	}
	ret := YouTubeDOM{
		root:     root,
		currNode: nil,
	}
	return &ret
}

//Parse get all the links
func (youTubeRoot *YouTubeDOM) Parse(root *html.Node, tagName, attrKey, attrVal string, targetNum int) {
	node := dfs(root, tagName, attrKey, attrVal, targetNum)
	youTubeRoot.currNode = node
}

func searchData(raw *[]byte, head int, data *[][]byte) int {
	cache := make([]byte, 0)
	for i := head; i < len(*raw); i++ {
		if (*raw)[i] == '{' {
			i = searchData(raw, i+1, data)
		} else if (*raw)[i] == '}' {
			cache = append([]byte{'{'}, cache...)
			cache = append(cache, '}')
			*data = append(*data, cache)
			return i
		} else {
			cache = append(cache, (*raw)[i])
		}
	}
	return len(*raw)
}

func getData(raw []byte) [][]byte {
	res := make([][]byte, 0)
	searchData(&raw, 0, &res)
	return res
}

func urlFormat(raw []byte) []byte {
	res := make([]byte, 0)
	for i := 0; i < len(raw); i++ {
		flag := true
		for _, kv := range URLFormater {
			if i+len(kv.key) < len(raw) && CompareSlice(raw[i:i+len(kv.key)], kv.key) {
				flag = false
				res = append(res, kv.val...)
				i += len(kv.key) - 1
			}
		}
		if flag {
			res = append(res, raw[i])
		}
	}
	return res
}

//Render return the html of a html node
func (youTubeRoot *YouTubeDOM) GetLinks() [](*VideoLink) {
	data := make([]byte, 0)
	w := writer{data: &data}
	html.Render(w, youTubeRoot.currNode)
	res := make([](*VideoLink), 0)
	for _, link := range getData(data) {
		l := VideoLink{}
		link = urlFormat(link)
		err := json.Unmarshal(link, &l)
		if err == nil && l.validate() {
			res = append(res, &l)
		}
	}
	return res
}
