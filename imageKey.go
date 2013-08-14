package main

import (
	"os"
)

// sort keys on file modified time in decreasing order (newer file first)
type keySorter struct {
	im *imageManager
}

func (s keySorter) Swap(i, j int) {
	s.im.keys[i], s.im.keys[j] = s.im.keys[j], s.im.keys[i]
}

func (s keySorter) Len() int {
	return len(s.im.keys)
}

func (s keySorter) Less(i, j int) bool {
	fsi, err := os.Stat(s.im.files[s.im.keys[i]])
	if err != nil {
		return false
	}
	fsj, err := os.Stat(s.im.files[s.im.keys[j]])
	if err != nil {
		return true
	}
	return fsi.ModTime().Unix() > fsj.ModTime().Unix()
}
