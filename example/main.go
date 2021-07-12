package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/davecgh/go-spew/spew"
	"github.com/ii64/protoc-gen-gohttpclient/example/gen"
)

var (
	debug = false
)

type transport struct {
	RoundTripHandler func(*http.Request) (*http.Response, error)
}

func (tr *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if f := tr.RoundTripHandler; f != nil {
		return f(req)
	}
	err = fmt.Errorf("roundtrip handler empty")
	return
}

func main() {
	hclient := &http.Client{}
	parentTransport := http.DefaultTransport
	hclient.Transport = &transport{
		RoundTripHandler: func(req *http.Request) (res *http.Response, err error) {
			if bb, err := httputil.DumpRequest(req, true); err == nil && debug {
				fmt.Printf("%s\n\n", bb)
			}
			if res, err = parentTransport.RoundTrip(req); err != nil {
				return
			}
			if bb, err := httputil.DumpResponse(res, true); err == nil && debug {
				fmt.Printf("%s\n\n", bb)
			}
			return
		},
	}

	client := gen.NewGreetServiceHTTPClient(
		"https://jsonplaceholder.typicode.com",
		hclient)
	fmt.Printf("Requesting GreetService.GetPost...\n")
	res, err := client.GetPost(context.Background(), &gen.GetPostRequest{
		Id: 7,
	})
	if err != nil {
		panic(err)
	}
	spew.Dump(res)

	fmt.Printf("\n\nRequesting KitsuService...\n")
	kcl := gen.NewKitsuServiceHTTPClient(
		"https://kitsu.io",
		hclient,
	)
	res2, err := kcl.GetAnime(context.Background(), &gen.KitsuAnimeRequest{
		Page: map[string]string{
			"limit":  "1",
			"offset": "0",
		},
	})
	if err != nil {
		panic(err)
	}
	spew.Dump(res2)

	fmt.Printf("Done\n")
}
