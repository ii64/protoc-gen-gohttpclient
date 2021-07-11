package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/ii64/protoc-gen-gohttpclient/example/gen"
)

func main() {
	client := gen.NewGreetServiceHTTPClient(
		"https://jsonplaceholder.typicode.com",
		&http.Client{})
	fmt.Printf("Requesting...\n")
	res, err := client.GetPost(context.Background(), &gen.GetPostRequest{
		Id: 7,
	})
	if err != nil {
		panic(err)
	}
	spew.Printf("%#+v\n", res)
	fmt.Printf("Done\n")
}
