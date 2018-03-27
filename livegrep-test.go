package main

import (
	"testing"
)

func TestQuery(t *testing.T) {
	l := &Livegrep{URL: "livegrep.com"}
	q := l.NewQuery("name")

	_, err := l.Query(q)
	if err != nil {
		// TODO: Make this more sophisticated
		panic(err)
	}

}
