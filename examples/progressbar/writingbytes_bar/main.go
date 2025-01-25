package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gngtwhh/gocui/pb"
)

func main() {
	req, _ := http.NewRequest("GET", "https://studygolang.com/dl/golang/go1.23.5.src.tar.gz", nil)
	req.Header.Add("Accept-Encoding", "identity")
	resp, err := http.DefaultClient.Do(req)
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

	fmt.Println("downloading...")
	bar, _ := pb.NewProgressBar("[%bar] %percent %bytes", pb.WithWriter(), pb.WithTotal(resp.ContentLength))
	barWriter, _ := bar.RunWithWriter()
	if _, err := io.Copy(io.MultiWriter(f, barWriter), resp.Body); err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print("\ndone")
}
