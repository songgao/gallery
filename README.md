# gallery

## What does it do?
`gallery` exposes a local gallery as a web gallery. It uses `lightbox2`(http://lokeshdhakar.com/projects/lightbox2/) to present images.

All images (`*.jpg`, `*.jpeg`, `*.png`, `*.gif`) in the folder are scaned
recursively and indexed. Images are indexed by SHA-1 of [absolute path of the
image file] and [modification time of the image file]. `gallery` watches the
image folder and re-indexes when the modification time of the folder changes

There are three sizes for each image: Thumbnail (width = 240px), Large (width = 1024), and Original. Thumbnail and Large are cached in `groupcache`(https://github.com/golang/groupcache) while Original is always loaded from hard drive. 8 MB and 64 MB are allocated for Thumbnails and Larges respectively.

## What does it look like?
![screenshots](https://raw.github.com/songgao/gallery/master/contrib/screenshots.jpg)

## Installation
```
go get -u github.com/songgao/gallery
```

## Usage
```
Usage of gallery:
  -image="": path to the folder that has images (supported formats: .jpg, .png, .gif)
  -laddr="localhost:7428": http listening address
```

```
gallery -image=/path/to/image/folder -laddr=localhost:7428
```
