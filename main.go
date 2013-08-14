package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	fImagePath string
	fLAddr     string
)

func init() {
	flag.StringVar(&fImagePath, "image", "", "path to the folder that has images (supported formats: .jpg, .png, .gif)")
	flag.StringVar(&fLAddr, "laddr", "localhost:7428", "http listening address")
}

func parseFlagOrPrintDefaults() (flagsOK bool) {
	flag.Parse()
	if fImagePath == "" {
		return false
	}

	return true
}

func main() {
	if !parseFlagOrPrintDefaults() {
		flag.PrintDefaults()
		return
	}
	im, err := initImageManager(fImagePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	mux, err := buildMux(im)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = http.ListenAndServe(fLAddr, mux)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func buildMux(im *imageManager) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/large/", func(w http.ResponseWriter, r *http.Request) {
		chunks := strings.Split(r.RequestURI, "/")
		key := chunks[len(chunks)-1]
		w.Header().Add("Content-Type", "image/jpeg")
		content, err := im.getLarge(key)
		if err != nil {
			http.NotFound(w, r)
		} else {
			http.ServeContent(w, r, key, time.Time{}, content)
		}
	})
	mux.HandleFunc("/thumbnail/", func(w http.ResponseWriter, r *http.Request) {
		chunks := strings.Split(r.RequestURI, "/")
		key := chunks[len(chunks)-1]
		w.Header().Add("Content-Type", "image/jpeg")
		content, err := im.getThumbnail(key)
		if err != nil {
			http.NotFound(w, r)
		} else {
			http.ServeContent(w, r, key, time.Time{}, content)
		}
	})
	mux.HandleFunc("/original/", func(w http.ResponseWriter, r *http.Request) {
		chunks := strings.Split(r.RequestURI, "/")
		key := chunks[len(chunks)-1]
		w.Header().Add("Content-Type", "image/jpeg")
		content, err := im.getOriginal(key)
		if err != nil {
			http.NotFound(w, r)
		} else {
			http.ServeContent(w, r, key, time.Time{}, content)
		}
	})

	rootPath, err := getRootPath()
	if err != nil {
		return nil, err
	}
	mux.Handle("/assets/", http.FileServer(http.Dir(rootPath)))

	tmpl, err := getTemplate("index.html")
	if err != nil {
		return nil, err
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		keys := im.getImageKeys()
		par := make([][]imageItem, 0)
		for i := 0; i < len(keys); i += 4 {
			rowCount := len(keys) - i
			if rowCount > 4 {
				rowCount = 4
			}
			row := make([]imageItem, rowCount)
			for j := 0; j < len(row); j++ {
				row[j].Original = fmt.Sprintf("/original/%s", keys[i+j])
				row[j].Large = fmt.Sprintf("/large/%s", keys[i+j])
				row[j].Thumbnail = fmt.Sprintf("/thumbnail/%s", keys[i+j])
				row[j].Title = im.getImageName(keys[i+j])
			}
			par = append(par, row)
		}
		tmpl.Execute(w, par)
	})

	return mux, nil
}

type imageItem struct {
	Original  string
	Large     string
	Thumbnail string
	Title     string
}
