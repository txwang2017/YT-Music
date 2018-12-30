package fortest

import (
	"os"
	"testing"
)

func TestT1(t *testing.T) {
	file, _ := os.Create("/Users/Tianxiong.wang/test/gozoo/fortest/ttt.txt")
	// file, err := os.OpenFile("/Users/Tianxiong.wang/test/gozoo/fortest/ttt.txt")
	file.WriteString("xxxasdfasdf")
	file.Close()
}
