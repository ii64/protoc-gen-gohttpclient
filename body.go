package main

// select which field to send in body or become url query

type bodySelector struct {
	field string
}

func bodyParse(body string) (*bodySelector, error) {
	return &bodySelector{
		field: body,
	}, nil
}
