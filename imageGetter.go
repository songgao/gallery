package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/groupcache"
	"github.com/nfnt/resize"
	"os"

	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
)

type widthLimitedImageGetter uint

func (g widthLimitedImageGetter) Get(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	files, ok := ctx.(map[string]string)
	if !ok {
		return errors.New("Context type wrong. (expecting map[string]string)")
	}
	filePath, ok := files[key]
	if !ok {
		return errors.New(fmt.Sprintf("Requested key (%s) does not exist.", key))
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(file)
	file.Close()
	if err != nil {
		return err
	}
	resized := resize.Resize(uint(g), 0, img, resize.NearestNeighbor)
	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, resized, nil)
	if err != nil {
		return err
	}
	return dest.SetBytes(buf.Bytes())
}
