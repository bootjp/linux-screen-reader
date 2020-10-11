package main

import (
	"os"

	reader "github.com/bootjp/google-tts-screenreader"
)

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/bootjp/bootjp-labs-5f1ff5ec3ac7.json")
	reader.NewObserve().Run()

}
