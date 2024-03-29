package main

import "google.golang.org/protobuf/compiler/protogen"

var (
	genVersion = "v1.1.10"
)

func main() {
	protogen.Options{}.Run(func(p *protogen.Plugin) error {
		for _, f := range p.Files {
			if f.Generate {
				if _, err := GenerateFile(p, f); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
