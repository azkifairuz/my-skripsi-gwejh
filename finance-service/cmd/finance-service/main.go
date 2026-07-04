package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("finance service is running")

	if os.Getenv("RUN_FOREVER") == "true" {
		for {
			time.Sleep(time.Hour)
		}
	}
}
