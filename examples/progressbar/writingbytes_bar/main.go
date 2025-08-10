package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gngtwhh/gocui/pb"
)

func main() {
	// The server only supports weak ciphers, which were disabled by default in Go 1.22.
	os.Setenv("GODEBUG", "tlsrsakex=1")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "https://studygolang.com/dl/golang/go1.23.5.src.tar.gz", nil)
	req.Header.Add("Accept-Encoding", "identity")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	f, _ := os.OpenFile("go1.23.5.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	defer func() {
		f.Close()
		if err := os.Remove("go1.23.5.src.tar.gz"); err != nil {
			panic(err)
		}
	}()

	fmt.Println("downloading...go1.23.5.src.tar.gz")
	// bar, _ := pb.NewProgressBar("[%bar] %percent %bytes", pb.WithWriter(), pb.WithTotal(resp.ContentLength))
	bar, _ := pb.NewProgressBar("[%bar] %percent %bytes", pb.WithWriter())
	barWriter, _ := bar.RunWithWriter(resp.ContentLength)
	if _, err := io.Copy(io.MultiWriter(f, barWriter), resp.Body); err != nil {
		fmt.Print(err.Error())
	}
	time.Sleep(time.Millisecond)
	fmt.Print("\ndone\n")
}
