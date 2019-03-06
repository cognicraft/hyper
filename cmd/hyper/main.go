package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/cognicraft/hyper"
)

func main() {
	flag.Parse()

	c := hyper.NewClient()
	item, err := c.Fetch(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	bs, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bs))
}
