package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("nlp service is running")

	if os.Getenv("RUN_FOREVER") == "true" {
		for {
			time.Sleep(time.Hour)
		}
	}
}
