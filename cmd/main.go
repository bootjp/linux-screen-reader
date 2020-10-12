package main

import (
	"fmt"
	"os"

	reader "github.com/bootjp/google-tts-screenreader"
)

func main() {
	if e := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); e == "" {
		fmt.Println("require environment GOOGLE_APPLICATION_CREDENTIALS")
		os.Exit(1)
	}

	reader.NewObserve().Run()
}
