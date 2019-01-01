package main

func main() {
	y := YouTubeAudio{}
	y.GetAudioMeta(`https://www.youtube.com/watch?v=VC2rAxRID9s`)
	y.Download("fengzhengwu.mp3")
}
