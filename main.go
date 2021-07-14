package main

import (
	"helmup/cmd"
	"io/ioutil"
	"log"
)

func main() {
	// Disable helm internal logging.
	log.SetOutput(ioutil.Discard)

	cmd.Execute()
}
