package main

import (
	"os"
	"path/filepath"
)

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
