package main

import (
	"github.com/gngtwhh/gocui/pb"
	"io"
	"net/http"
	"os"
)

func main() {
	req, _ := http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile("go1.14.2.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar, _ := pb.NewProgressBar("[%bar] %percent %rate %bytes", pb.WithWriter(), pb.WithTotal(resp.ContentLength))
	barWriter, stop := bar.RunWithWriter()
	if _, err := io.Copy(io.MultiWriter(f, barWriter), resp.Body); err != nil {
		close(stop)
		panic(err)
	}
}
