package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/golang/groupcache"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var acceptedImageExt = []string{".png", ".gif", ".jpg", ".jpeg"}

type imageManager struct {
	imageFolderPath string

	files         map[string]string // SHA-1 -> filepath
	files_reverse map[string]string // filepath -> SHA-1
	keys          []string
	mu            *sync.RWMutex

	large     *groupcache.Group
	thumbnail *groupcache.Group
}

func initImageManager(imageFolderPath string) (*imageManager, error) {
	i := &imageManager{imageFolderPath: imageFolderPath}
	i.mu = new(sync.RWMutex)
	i.initGroupCache()
	i.scanFiles()
	err := i.startWatchingFS()
	if err != nil {
		return nil, err
	}
	return i, err
}

func (i *imageManager) initGroupCache() {
	i.large = groupcache.NewGroup("large", 128*1024*1024, widthLimitedImageGetter(1024))
	i.thumbnail = groupcache.NewGroup("thumbnail", 32*1024*1024, widthLimitedImageGetter(240))
}

func (i *imageManager) getImageKeys() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.keys
}

func (i *imageManager) getImageName(key string) string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return filepath.Base(i.files[key])
}

func (i *imageManager) getLarge(key string) (io.ReadSeeker, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	bv := groupcache.ByteView{}
	err := i.large.Get(i.files, key, groupcache.ByteViewSink(&bv))
	return bv.Reader(), err
}

func (i *imageManager) getThumbnail(key string) (io.ReadSeeker, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	bv := groupcache.ByteView{}
	err := i.thumbnail.Get(i.files, key, groupcache.ByteViewSink(&bv))
	return bv.Reader(), err
}

func (i *imageManager) getOriginal(key string) (io.ReadSeeker, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	f, err := os.Open(i.files[key])
	return f, err
}

func (i *imageManager) scanFiles() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.files = make(map[string]string)
	i.files_reverse = make(map[string]string)
	i.keys = make([]string, 0)

	imageHasher := sha1.New()
	filepath.Walk(i.imageFolderPath, func(filePath string, info os.FileInfo, err error) error {
		if err == nil {
			fstat, err := os.Stat(filePath)
			if err == nil && fstat.Mode().IsRegular() && in(acceptedImageExt, strings.ToLower(path.Ext(filePath))) {
				absPath, err := filepath.Abs(filePath)
				if err == nil {
					io.WriteString(imageHasher, absPath)
					io.WriteString(imageHasher, fstat.ModTime().String())
					sha1 := fmt.Sprintf("%x", imageHasher.Sum(nil))
					i.files[sha1] = absPath
					i.files_reverse[absPath] = sha1
					i.keys = append(i.keys, sha1)
					imageHasher.Reset()
				}
			}
		}
		return nil
	})
	sort.Sort(keySorter{i})
}

func (i *imageManager) startWatchingFS() error {
	fs, err := os.Stat(i.imageFolderPath)
	if err != nil {
		return err
	}
	lastMod := fs.ModTime().Unix()
	go func() {
		for {
			time.Sleep(4 * time.Second)
			fs, err = os.Stat(i.imageFolderPath)
			if err == nil {
				mod := fs.ModTime().Unix()
				if lastMod != mod {
					i.scanFiles()
					lastMod = mod
				}
			}
		}
	}()
	return nil
}
