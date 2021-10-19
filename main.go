package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/kurin/blazer/b2"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed english.txt
var english string

//go:embed index.html
var indexHtml string
var wordlists []string
var wordlistsCount int

func main() {
	baseUrl := os.Getenv("BASE_URL")
	baseDownloadUrl := os.Getenv("BASE_DOWNLOAD_URL")

	ctx := context.Background()
	b2, err := b2.NewClient(ctx, os.Getenv("B2_KEY"), os.Getenv("B2_SECRET"))
	if err != nil {
		log.Fatalln(err)
	}
	bucket, err := b2.Bucket(ctx, os.Getenv("B2_BUCKET"))
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprintf(w, indexHtml, baseUrl)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			return
		case "POST":
			key := randomKey(3, "-")
			obj := bucket.Object(key)
			objWriter := obj.NewWriter(ctx)
			defer objWriter.Close()
			_, err := io.Copy(objWriter, r.Body)
			if err != nil {
				log.Printf("Error uploading trash, %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "%s%s", baseDownloadUrl, key)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	wordlists = strings.Split(english, "\n")
	wordlistsCount = len(wordlists)
	rand.Seed(time.Now().UnixNano())
}

func randomKey(keysize int, connector string) string {
	var t []string
	for i := 0; i < keysize; i++ {
		t = append(t, wordlists[rand.Intn(wordlistsCount)])
	}
	return strings.Join(t, connector)
}
