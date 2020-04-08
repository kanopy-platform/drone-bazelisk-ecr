package main

import (
	"log"
)

func main() {
	// get new plugin object
	p := newPlugin()

	// run bazelisk
  err := p.run()
	if err != nil {
		log.Fatal(err)
	}
}
