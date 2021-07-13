package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/ii64/protoc-gen-gohttpclient/example/gen"
	lib "github.com/ii64/protoc-gen-gohttpclient/lib"
)

var (
	debug = true
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

	kcl := gen.NewKitsuServiceHTTPClient(
		"https://kitsu.io",
		hclient,
		lib.WithIgnoreUnknownField,
		lib.WithHTTPResponseValidator(func(res *http.Response) error {
			if res.StatusCode != 200 {
				return fmt.Errorf("status is not 200")
			}
			return nil
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/s/kitsu/getAnime", func(rw http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		respQueryString := r.URL.Query().Get("qstr") != ""
		if id == "" {
			rw.WriteHeader(400)
			return
		}

		resp, err := kcl.GetAnime(r.Context(), &gen.KitsuAnimeRequest{
			Id: id,
		})
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
			return
		}
		resp.ContentSource = 2
		rw.Header().Set("Content-Type", "application/json")
		var b []byte
		if respQueryString {
			b = []byte(resp.QueryString().Encode())
		} else {
			mo := protojson.MarshalOptions{}
			if b, err = mo.Marshal(resp); err != nil {
				rw.WriteHeader(500)
				return
			}
		}
		rw.Write(b)
	})

	mux.HandleFunc("/s/kitsu/getAnimeList", func(rw http.ResponseWriter, r *http.Request) {
		var (
			rLimit  string = "5"
			rOffset string = "0"
		)
		if limit := r.URL.Query().Get("limit"); limit != "" {
			rLimit = limit
		}
		if offset := r.URL.Query().Get("offset"); offset != "" {
			rOffset = offset
		}
		respQueryString := r.URL.Query().Get("qstr") != ""

		resp, err := kcl.GetAnimeList(r.Context(), &gen.KitsuAnimeListRequest{
			Page: map[string]string{
				"limit":  rLimit,
				"offset": rOffset,
			},
		})
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
			return
		}
		resp.ContentSource = 2
		rw.Header().Set("Content-Type", "application/json")
		var b []byte
		if respQueryString {
			b = []byte(resp.QueryString().Encode())
		} else {
			mo := protojson.MarshalOptions{}
			if b, err = mo.Marshal(resp); err != nil {
				rw.WriteHeader(500)
				return
			}
		}
		rw.Write(b)
	})
	addr := ":8888"
	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	fmt.Printf("HTTP Server listening on %s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
